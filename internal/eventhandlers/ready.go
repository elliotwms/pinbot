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
		err := s.UpdateGameStatus(0, build.Version)
		if err != nil {
			log.WithError(err).Error("Could not update game status")
		}

		// ensure the Pin command is created
		if _, err := s.ApplicationCommandCreate(config.ApplicationID, "", commands.Pin); err != nil {
			log.WithError(err).Error("Could not register import command")
		}
	}
}
