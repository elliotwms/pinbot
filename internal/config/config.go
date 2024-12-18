package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/build"
	"github.com/sirupsen/logrus"
)

const DefaultIntents = discordgo.IntentsGuilds |
	discordgo.IntentsGuildMessages |
	discordgo.IntentsGuildMessageReactions

const DefaultPermissions = discordgo.PermissionViewChannel |
	discordgo.PermissionSendMessages |
	discordgo.PermissionAddReactions

var (
	Token           string
	ApplicationID   string
	HealthCheckAddr string
	LogLevel        logrus.Level
	Intents         discordgo.Intent
	Permissions     int
)

var once sync.Once

func Configure() {
	once.Do(func() {
		Token = mustGetEnv("TOKEN")
		ApplicationID = mustGetEnv("APPLICATION_ID")
		HealthCheckAddr = os.Getenv("HEALTH_CHECK_ADDR")

		l, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
		if err != nil {
			LogLevel = logrus.InfoLevel
		} else {
			LogLevel = l
		}

		Intents = DefaultIntents
		if s := os.Getenv("INTENTS"); s != "" {
			if i, err := strconv.Atoi(s); err == nil {
				Intents = discordgo.Intent(i)
			}
		}

		Permissions = DefaultPermissions
		if s := os.Getenv("PERMISSIONS"); s != "" {
			if i, err := strconv.Atoi(s); err == nil {
				Permissions = i
			}
		}
	})
}

func Output(showSensitive bool) logrus.Fields {
	fields := logrus.Fields{
		"APPLICATION_ID":    ApplicationID,
		"HEALTH_CHECK_ADDR": HealthCheckAddr,
		"LOG_LEVEL":         LogLevel,
		"INTENTS":           Intents,
		"PERMISSIONS":       Permissions,
		"install_url":       BuildInstallURL().String(),
		"version":           build.Version,
	}

	if showSensitive {
		fields["TOKEN"] = Token
	}

	return fields
}

func mustGetEnv(s string) string {
	token := os.Getenv(s)
	if token == "" {
		panic(fmt.Sprintf("Missing '%s'", s))
	}
	return token
}

func BuildInstallURL() *url.URL {
	u, _ := url.Parse("https://discord.com/oauth2/authorize")

	q := u.Query()
	q.Add("client_id", ApplicationID)
	q.Add("permissions", fmt.Sprintf("%d", Permissions))
	q.Add("scope", "applications.commands bot")
	u.RawQuery = q.Encode()

	return u
}
