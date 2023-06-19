package commandhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/commands"
	"github.com/sirupsen/logrus"
)

type RegisterCommandsCommand struct {
	ApplicationID, GuildID string
}

func RegisterCommandsCommandHandler(c *RegisterCommandsCommand, s *discordgo.Session, log *logrus.Entry) {
	_, err := s.ApplicationCommandCreate(c.ApplicationID, c.GuildID, commands.Import)
	if err != nil {
		log.WithField("guild_id", c.GuildID).WithError(err).Error("Could not register import command")
	}
}
