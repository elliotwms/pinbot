package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/pinbot"
)

func main() {
	// build the logger and session
	logLevel := getLogLevel(os.Getenv("LOG_LEVEL"))
	slog.SetLogLoggerLevel(logLevel)
	s := buildSession(logLevel)

	// configure the bot
	c := pinbot.NewConfig(s, mustGetEnv("APPLICATION_ID"))
	c.HealthCheckAddr = os.Getenv("HEALTH_CHECK_ADDR")
	c.GuildID = os.Getenv("GUILD_ID")
	c.Logger = slog.Default()

	// listen for signals
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// run the bot
	if err := pinbot.Run(c, ctx); err != nil {
		slog.Error("completed with error", "error", err)
		os.Exit(1)
	}
}

func getLogLevel(s string) (l slog.Level) {
	if s == "" {
		return slog.LevelInfo
	}

	if err := l.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}

	return l
}

func buildSession(l slog.Level) *discordgo.Session {
	s, err := discordgo.New("Bot " + mustGetEnv("TOKEN"))
	if err != nil {
		panic(err)
	}

	// set the discord log level
	if l <= slog.LevelDebug {
		s.LogLevel = discordgo.LogDebug
	}
	return s
}

func mustGetEnv(s string) string {
	token := os.Getenv(s)
	if token == "" {
		panic(fmt.Sprintf("Missing '%s'", s))
	}
	return token
}
