package commands

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

const (
	emojiPinned     = "📌"
	pinMessageColor = 0xbb0303
)

var Pin = &discordgo.ApplicationCommand{
	Name: "Pin",
	Type: discordgo.MessageApplicationCommand,
}

func PinMessageCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) (err error) {
	ctx := context.TODO()

	m := data.Resolved.Messages[data.TargetID]
	m.GuildID = i.GuildID // guildID is missing from message in resolved context

	log := slog.With("guild_id", i.GuildID, "channel_id", i.ChannelID, "message_id", m.ID)

	log.Debug("Starting pin message")

	pinned, err := isAlreadyPinned(ctx, s, i, m)
	if err != nil {
		log.Error("Could not check if message is already pinned", "error", err)
		return respond(ctx, s, i.Interaction, "💩 Temporary error, please retry")
	}
	if pinned {
		return respond(ctx, s, i.Interaction, "🔄 Message already pinned")
	}

	// get the source channel
	sourceChannel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Error("Could not determine source channel", "error", err)
		return respond(ctx, s, i.Interaction, "💩 Temporary error, please retry")
	}

	// determine the target pin channel for the message
	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		log.Error("Could not get guild", "error", err)
		return respond(ctx, s, i.Interaction, "💩 Temporary error, please retry")
	}

	targetChannel, err := getTargetChannel(g.Channels, sourceChannel)
	if err != nil {
		log.Error("Could not determine target channel", "error", err)
		return respond(ctx, s, i.Interaction, "💩 Temporary error, please retry")
	}
	log = log.With("target_channel_id", targetChannel.ID)

	// build the rich embed pin message
	forward := &discordgo.MessageSend{
		Reference: m.Forward(),
	}

	log.Debug("Sending pin message")
	pin, err := s.ChannelMessageSendComplex(targetChannel.ID, forward, discordgo.WithContext(ctx))
	if err != nil {
		log.Error("Could not send pin message", "error", err)
		return respond(ctx, s, i.Interaction, "🙅 Could not send pin message. Please ensure bot has permission to post in "+targetChannel.Mention())
	}

	// build a secondary info message to accompany the forward
	infoMessage := buildPinInfoMessage(sourceChannel, m, pin, i.Member.User)

	log.Debug("Sending pin info message")
	_, err = s.ChannelMessageSendComplex(targetChannel.ID, infoMessage, discordgo.WithContext(ctx))
	if err != nil {
		log.Error("Could not send info message", "error", err)
		return respond(ctx, s, i.Interaction, "🙅 Could not send pin message. Please ensure bot has permission to post in "+targetChannel.Mention())
	}

	// mark the message as done
	if err := s.MessageReactionAdd(m.ChannelID, m.ID, emojiPinned, discordgo.WithContext(ctx)); err != nil {
		log.Error("Could not react to message", "error", err)
	}

	log.Info("Pinned message", "pin_message_id", pin.ID)

	return respond(ctx, s, i.Interaction, "📌 Pinned: "+url(i.GuildID, pin.ChannelID, pin.ID))
}

func respond(ctx context.Context, s *discordgo.Session, i *discordgo.Interaction, c string) error {
	_, err := s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Content: &c,
	}, discordgo.WithContext(ctx))

	return err
}

func url(guildID, channelID, messageID string) string {
	return fmt.Sprintf(
		"https://discord.com/channels/%s/%s/%s",
		guildID,
		channelID,
		messageID,
	)
}

func buildPinInfoMessage(sourceChannel *discordgo.Channel, m *discordgo.Message, pin *discordgo.Message, pinnedBy *discordgo.User) *discordgo.MessageSend {
	u := url(sourceChannel.GuildID, pin.ChannelID, pin.ID)
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    m.Author.Username,
			IconURL: m.Author.AvatarURL(""),
			URL:     u,
		},
		Title: "📌 Pinned",
		Color: pinMessageColor,
		URL:   u,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Channel",
				Value:  sourceChannel.Mention(),
				Inline: true,
			},
		},
	}

	if pinnedBy != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Pinned by",
			Value:  pinnedBy.Mention(),
			Inline: true,
		})
	}

	return &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
	}
}

func isAlreadyPinned(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, m *discordgo.Message) (bool, error) {
	acks, err := s.MessageReactions(m.ChannelID, m.ID, emojiPinned, 0, "", "", discordgo.WithContext(ctx))
	if err != nil {
		return false, err
	}

	for _, ack := range acks {
		if ack.ID == i.AppID {
			return true, nil
		}
	}

	return false, nil
}

// getTargetChannel returns the target pin channel for a given channel #channel in the following order:
// #channel-pins (a specific pin channel)
// #pins (a generic pin channel)
// #channel (the channel itself)
func getTargetChannel(channels []*discordgo.Channel, origin *discordgo.Channel) (*discordgo.Channel, error) {
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
