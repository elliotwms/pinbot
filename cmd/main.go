package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/bot"
	"github.com/elliotwms/pinbot/internal/config"
	"github.com/elliotwms/pinbot/internal/pinbot"
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

	b := bot.New(config.ApplicationID, s, log)

	//if config.TestGuildID != "" {
	//	b.WithTestGuildID(config.TestGuildID)
	//}

	if config.HealthCheckAddr != "" {
		b.WithHealthCheck(config.HealthCheckAddr)
	}

	done := pinbot.RegisterHandlers(s, logrus.NewEntry(log))
	defer done()

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	if err := b.Run(ctx); err != nil {
		os.Exit(1)
	}
}
