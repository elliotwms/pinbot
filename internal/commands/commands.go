package commands

import "github.com/bwmarrin/discordgo"

const OptionChannel = "channel"

var Import = &discordgo.ApplicationCommand{
	Name:        "import",
	Description: "Import existing channel pins",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionChannel,
			Name:        OptionChannel,
			Description: "Channel",
			ChannelTypes: []discordgo.ChannelType{
				discordgo.ChannelTypeGuildText,
			},
			Required: false,
		},
	},
}
