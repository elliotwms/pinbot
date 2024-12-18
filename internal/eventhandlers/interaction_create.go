package eventhandlers

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/elliotwms/pinbot/internal/commands"
	"github.com/sirupsen/logrus"
)

func InteractionCreate(log *logrus.Entry) func(s *discordgo.Session, e *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, e *discordgo.InteractionCreate) {
		// only process application commands
		if e.Type != discordgo.InteractionApplicationCommand {
			return
		}

		err := s.InteractionRespond(e.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.WithError(err).Error("Failed to respond to InteractionCreate")
			return
		}

		command := e.ApplicationCommandData()
		switch command.Name {
		case commands.Pin.Name:
			err = commandhandlers.PinMessageCommandHandler(context.Background(), s, e, command)
		}
	}
}
