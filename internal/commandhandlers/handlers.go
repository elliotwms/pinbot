package commandhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/commands"
	"github.com/sirupsen/logrus"
)

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData, log *logrus.Entry) (string, error)

var Handlers = map[string]CommandHandler{
	commands.Pin.Name: pinMessageCommandHandler,
}
