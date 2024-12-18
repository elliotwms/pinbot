package commands

import "github.com/bwmarrin/discordgo"

var Pin = &discordgo.ApplicationCommand{
	Name: "Pin",
	Type: discordgo.MessageApplicationCommand,
}
