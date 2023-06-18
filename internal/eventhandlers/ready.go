package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func Ready(log *logrus.Entry) func(s *discordgo.Session, _ *discordgo.Ready) {
	return func(s *discordgo.Session, _ *discordgo.Ready) {
		log.Info("I am ready for action")
		err := s.UpdateGameStatus(0, "Transcribing pins")
		if err != nil {
			log.WithError(err).Error("Could not update game status")
			return
		}
	}
}
