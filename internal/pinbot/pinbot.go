package pinbot

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/bot"
	"github.com/elliotwms/bot/interactions"
	"github.com/elliotwms/pinbot/internal/commands"
	"github.com/elliotwms/pinbot/internal/eventhandlers"
)

const intents = discordgo.IntentsGuilds

type Config struct {
	Session         *discordgo.Session
	ApplicationID   string
	HealthCheckAddr string
	GuildID         string
	Logger          *slog.Logger
}

func NewConfig(s *discordgo.Session, appID string) Config {
	return Config{
		Session:       s,
		ApplicationID: appID,
	}
}

// Run builds and runs the bot.
func Run(config Config, ctx context.Context) error {
	r := interactions.NewRouter(
		interactions.WithDeferredResponse(true),
		interactions.WithLogger(slog.Default()),
	)

	b := bot.
		New(config.ApplicationID, config.Session).
		WithLogger(config.Logger).
		WithRouter(r).
		WithIntents(intents).
		WithHandler(eventhandlers.Ready).
		WithMigrationEnabled(true).
		WithApplicationCommand(commands.Pin, commands.PinMessageCommandHandler)

	if config.HealthCheckAddr != "" {
		b.WithHealthCheck(config.HealthCheckAddr)
	}

	if config.GuildID != "" {
		b.WithGuildID(config.GuildID)
	}

	return b.Build().Run(ctx)
}
