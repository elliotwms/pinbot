package commands

import "github.com/bwmarrin/discordgo"

var Pin = &discordgo.ApplicationCommand{
	Name: "Pin",
	Type: discordgo.MessageApplicationCommand,
	IntegrationTypes: &[]discordgo.ApplicationIntegrationType{
		discordgo.ApplicationIntegrationGuildInstall,
	},
	Contexts: &[]discordgo.InteractionContextType{
		discordgo.InteractionContextGuild,
	},
}
