package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func MessageReactionAdd(log *logrus.Entry) func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	return func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
		log.WithField("emoji", e.Emoji.Name).Debug("Received reaction")

		if e.Emoji.Name != "📌" {
			// only react to pin emojis
			return
		}

		// todo notify channel of new pin behaviour?
	}
}
