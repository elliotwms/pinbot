package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/build"
	"github.com/elliotwms/pinbot/internal/commands"
	"github.com/elliotwms/pinbot/internal/config"
	"github.com/sirupsen/logrus"
)

func Ready(log *logrus.Entry) func(s *discordgo.Session, _ *discordgo.Ready) {
	return func(s *discordgo.Session, _ *discordgo.Ready) {
		log.Info("I am ready for action")
		err := s.UpdateGameStatus(0, build.Version)
		if err != nil {
			log.WithError(err).Error("Could not update game status")
			return
		}

		// check if Pin command exists, create if not
		cs, err := s.ApplicationCommands(config.ApplicationID, "")
		if err != nil {
			log.WithError(err).Error("Could not get application commands")
			return
		}

		for _, c := range cs {
			if c.Name == commands.Pin.Name && c.Type == commands.Pin.Type {
				log.Info("Pin command already exists")
				return
			}
		}

		log.Info("Creating Pin command")
		_, err = s.ApplicationCommandCreate(config.ApplicationID, "", commands.Pin)
		if err != nil {
			log.WithError(err).Error("Could not create command")
		}
	}
}
