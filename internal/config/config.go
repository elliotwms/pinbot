package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"

	"github.com/bwmarrin/discordgo"
)

const DefaultIntents = discordgo.IntentsGuilds |
	discordgo.IntentsGuildMessages |
	discordgo.IntentsGuildMessageReactions

var (
	Token           string
	ApplicationID   string
	HealthCheckAddr string
	LogLevel        slog.Level
	Intents         discordgo.Intent
)

var once sync.Once

func Configure() {
	once.Do(func() {
		Token = mustGetEnv("TOKEN")
		ApplicationID = mustGetEnv("APPLICATION_ID")
		HealthCheckAddr = os.Getenv("HEALTH_CHECK_ADDR")

		level := os.Getenv("LOG_LEVEL")
		if level == "" {
			LogLevel = slog.LevelInfo
		} else {
			err := LogLevel.UnmarshalText([]byte(level))
			if err != nil {
				panic(err)
			}
		}

		Intents = DefaultIntents
		if s := os.Getenv("INTENTS"); s != "" {
			if i, err := strconv.Atoi(s); err == nil {
				Intents = discordgo.Intent(i)
			}
		}
	})
}

func mustGetEnv(s string) string {
	token := os.Getenv(s)
	if token == "" {
		panic(fmt.Sprintf("Missing '%s'", s))
	}
	return token
}
