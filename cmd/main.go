package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/bot"
	"github.com/elliotwms/pinbot/internal/commands"
	"github.com/elliotwms/pinbot/internal/config"
	"github.com/elliotwms/pinbot/internal/eventhandlers"
)

func main() {
	config.Configure()

	slog.SetLogLoggerLevel(config.LogLevel)

	s, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		panic(err)
	}

	// set the discord log level
	if config.LogLevel == slog.LevelDebug {
		s.LogLevel = discordgo.LogDebug
	}

	b := bot.
		New(config.ApplicationID, s).
		WithLogger(slog.Default()).
		WithIntents(config.Intents).
		WithHandler(eventhandlers.Ready).
		WithMigrationEnabled(true).
		WithApplicationCommand(commands.Pin, commands.PinMessageCommandHandler)

	if config.HealthCheckAddr != "" {
		b.WithHealthCheck(config.HealthCheckAddr)
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	if err := b.Build().Run(ctx); err != nil {
		os.Exit(1)
	}
}
