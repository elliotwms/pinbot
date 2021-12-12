package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/pinbot"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("Missing TOKEN")
	}

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	s.LogLevel = discordgo.LogDebug

	bot := pinbot.New(s, log)

	if id := os.Getenv("TEST_GUILD_ID"); id != "" {
		bot.WithTestGuildID(id)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	if err := bot.StartSession(sc); err != nil {
		os.Exit(1)
	}
}
