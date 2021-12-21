package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/sirupsen/logrus"
)

func MessageReactionAdd(log *logrus.Entry, testGuildID string) func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	return func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
		log.WithField("emoji", e.Emoji.Name).Info("Received reaction")

		if e.Emoji.Name != "📌" {
			// only react to pin emojis
			return
		}

		if testGuildID != "" && testGuildID != e.GuildID {
			log.Info("Skipping non-test guild")
			return
		}

		m, err := s.ChannelMessage(e.ChannelID, e.MessageID)
		if err != nil {
			log.WithError(err).Error("Could not get channel message")
			return
		}

		commandhandlers.PinMessageCommandHandler(&commandhandlers.PinMessageCommand{
			GuildID:  e.GuildID,
			Message:  m,
			PinnedBy: e.Member.User,
		}, s, log)
	}
}
