package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/config"
	"github.com/sirupsen/logrus"
)

func GuildCreate(log *logrus.Entry) func(s *discordgo.Session, e *discordgo.GuildCreate) {
	return func(s *discordgo.Session, e *discordgo.GuildCreate) {
		log.Debug("Guild info received:", e.Name)

		err := cleanupOldCommands(config.ApplicationID, e.Guild.ID, s)
		if err != nil {
			log.WithError(err).Error("Error cleaning up old commands")
		}
	}
}

func cleanupOldCommands(applicationID, guildID string, s *discordgo.Session) error {
	cmds, err := s.ApplicationCommands(applicationID, guildID)
	if err != nil {
		return err
	}

	for _, cmd := range cmds {
		// remove the old import command
		if cmd.Name != "import" {
			err = s.ApplicationCommandDelete(applicationID, guildID, cmd.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
