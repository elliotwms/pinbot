package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/bot"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/elliotwms/pinbot/internal/commands"
	"github.com/elliotwms/pinbot/internal/config"
	"github.com/elliotwms/pinbot/internal/eventhandlers"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	config.Configure()
	log.SetLevel(config.LogLevel)

	s, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		panic(err)
	}

	if log.IsLevelEnabled(logrus.TraceLevel) {
		s.LogLevel = discordgo.LogDebug
	}

	b := bot.
		New(config.ApplicationID, s).
		WithIntents(config.Intents).
		WithHandlers(eventhandlers.List(logrus.NewEntry(log))).
		WithMigrationEnabled(true).
		WithApplicationCommand(commands.Pin, commandhandlers.PinMessageCommandHandler)

	if config.HealthCheckAddr != "" {
		b.WithHealthCheck(config.HealthCheckAddr)
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	if err := b.Build().Run(ctx); err != nil {
		os.Exit(1)
	}
}
