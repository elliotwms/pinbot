package commandhandlers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const (
	emojiPinned = "📌"
)

const pinMessageColor = 0xbb0303

func PinMessageCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) (userFeedback string, err error) {
	m := data.Resolved.Messages[data.TargetID]
	m.GuildID = i.GuildID // guildID is missing from message in resolved context

	log := logrus.WithFields(logrus.Fields{
		"guild_id":   i.GuildID,
		"channel_id": i.ChannelID,
		"message_id": m.ID,
	})

	log.Debug("Pinning message")

	pinned, err := isAlreadyPinned(s, m)
	if err != nil {
		// proceed but assume the message is not already pinned
		log.WithError(err).Error("Could not check if message is already pinned. Assuming unpinned...")
	}
	if pinned {
		return "🔄 Message already pinned", nil
	}

	sourceChannel, err := s.Channel(m.ChannelID)
	if err != nil {
		return "💩 Temporary error, please retry", fmt.Errorf("determine source channel: %w", err)
	}

	// determine the target pin channel for the message
	targetChannel, err := getTargetChannel(s, i.GuildID, sourceChannel)
	if err != nil {
		return "💩 Temporary error, please retry", fmt.Errorf("determine target channel: %w", err)
	}

	l := log.WithField("target_channel_id", targetChannel.ID)

	// build the rich embed pin message
	pinMessage := buildPinMessage(sourceChannel, m, i.Member.User)

	// send the pin message
	pin, err := s.ChannelMessageSendComplex(targetChannel.ID, pinMessage)
	if err != nil {
		return "🙅 Could not send pin message. Please check bot permissions", fmt.Errorf("send pin message: %w", err)
	}

	// mark the message as done
	react(s, m, emojiPinned, l)

	l.Info("Pinned message")

	return "📌 Pinned: " + url(i.GuildID, pin.ChannelID, pin.ID), nil
}

func react(s *discordgo.Session, m *discordgo.Message, emoji string, l *logrus.Entry) {
	if err := s.MessageReactionAdd(m.ChannelID, m.ID, emoji); err != nil {
		l.WithError(err).Error("Could not react to message")
	}
}

func url(guildID, channelID, messageID string) string {
	return fmt.Sprintf(
		"https://discord.com/channels/%s/%s/%s",
		guildID,
		channelID,
		messageID,
	)
}

func buildPinMessage(sourceChannel *discordgo.Channel, m *discordgo.Message, pinnedBy *discordgo.User) *discordgo.MessageSend {
	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "Channel",
			Value:  sourceChannel.Mention(),
			Inline: true,
		},
	}

	u := url(sourceChannel.GuildID, m.ChannelID, m.ID)
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    m.Author.Username,
			IconURL: m.Author.AvatarURL(""),
			URL:     u,
		},
		Title:       "📌 Pinned",
		Color:       pinMessageColor,
		Description: m.Content,
		URL:         u,
		Timestamp:   m.Timestamp.Format(time.RFC3339),
	}

	if pinnedBy != nil {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Pinned by",
			Value:  pinnedBy.Mention(),
			Inline: true,
		})
	}

	embed.Fields = fields

	pinMessage := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
	}

	// If there are multiple attachments then add them to separate embeds
	for i, a := range m.Attachments {
		if a.Width == 0 || a.Height == 0 {
			// only embed images
			continue
		}
		e := &discordgo.MessageEmbedImage{URL: a.URL}

		if i == 0 {
			// add the first image to the existing embed
			pinMessage.Embeds[0].Image = e
		} else {
			// add any other images to their own embed
			pinMessage.Embeds = append(pinMessage.Embeds, &discordgo.MessageEmbed{
				Type:  discordgo.EmbedTypeImage,
				Color: pinMessageColor,
				Image: e,
			})
		}
	}

	// preserve the existing embeds
	pinMessage.Embeds = append(pinMessage.Embeds, m.Embeds...)

	return pinMessage
}

func isAlreadyPinned(s *discordgo.Session, m *discordgo.Message) (bool, error) {
	acks, err := s.MessageReactions(m.ChannelID, m.ID, emojiPinned, 0, "", "")
	if err != nil {
		return false, err
	}

	for _, ack := range acks {
		if ack.ID == s.State.User.ID {
			return true, nil
		}
	}

	return false, nil
}

// getTargetChannel returns the target pin channel for a given channel #channel in the following order:
// #channel-pins (a specific pin channel)
// #pins (a generic pin channel)
// #channel (the channel itself)
func getTargetChannel(s *discordgo.Session, guildID string, origin *discordgo.Channel) (*discordgo.Channel, error) {
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return nil, err
	}

	// use the same channel by default
	channel := origin

	// check for #channel-pins first
	for _, c := range channels {
		if c.Name == channel.Name+"-pins" && c.Type == discordgo.ChannelTypeGuildText {
			return c, nil
		}
	}

	// fallback to general pins channel
	for _, c := range channels {
		if c.Name == "pins" && c.Type == discordgo.ChannelTypeGuildText {
			return c, nil
		}
	}

	return channel, nil
}
