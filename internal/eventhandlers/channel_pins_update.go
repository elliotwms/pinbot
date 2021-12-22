package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/elliotwms/pinbot/internal/config"
	"github.com/sirupsen/logrus"
)

func ChannelPinsUpdate(log *logrus.Entry) func(s *discordgo.Session, e *discordgo.ChannelPinsUpdate) {
	return func(s *discordgo.Session, e *discordgo.ChannelPinsUpdate) {
		if !config.ShouldActOnGuild(e.GuildID) {
			return
		}

		commandhandlers.ImportChannelCommandHandler(&commandhandlers.ImportChannelCommand{
			ChannelID: e.ChannelID,
			GuildID:   e.GuildID,
		}, s, log)
	}
}
