package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/sirupsen/logrus"
)

func InteractionCreate(log *logrus.Entry) func(s *discordgo.Session, e *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		data := i.ApplicationCommandData()

		l := log.WithFields(map[string]interface{}{
			"guild_id":       i.GuildID,
			"channel_id":     i.ChannelID,
			"interaction_id": i.ID,
			"message_id":     data.TargetID,
			"interaction":    data.Name,
		})

		h, ok := commandhandlers.Handlers[data.Name]

		if !ok {
			l.Warn("Unexpected interaction")
			return
		}

		res, err := h(s, i, data, l)

		if err != nil {
			l.WithError(err).Error("Error handling interaction")
		}

		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: res,
			},
		})
	}
}
