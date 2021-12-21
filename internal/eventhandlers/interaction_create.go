package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/elliotwms/pinbot/internal/pinbot/commands"
	"github.com/sirupsen/logrus"
)

func InteractionCreate(log *logrus.Entry) func(s *discordgo.Session, e *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, e *discordgo.InteractionCreate) {
		if e.Type != discordgo.InteractionApplicationCommand {
			return
		}

		command := e.ApplicationCommandData()
		switch command.Name {
		case commands.Import.Name:
			channel := e.ChannelID
			for _, option := range command.Options {
				if option.Name == commands.OptionChannel {
					if c, ok := option.Value.(string); ok {
						channel = c
					}
				}
			}

			commandhandlers.ImportChannelCommandHandler(&commandhandlers.ImportChannelCommand{
				GuildID:   e.GuildID,
				ChannelID: channel,
			}, s, log)

			_ = s.InteractionRespond(e.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Starting import",
					Flags:   1 << 6,
				},
			})
		}
	}
}
