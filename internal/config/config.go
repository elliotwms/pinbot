package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
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
	Token            string
	ApplicationID    string
	HealthCheckAddr  string
	LogLevel         logrus.Level
	SelfPinEnabled   bool
	ExcludedChannels []string
	Intents          discordgo.Intent
	Permissions      int
)

var once sync.Once

func Configure() {
	once.Do(func() {
		Token = mustGetEnv("TOKEN")
		ApplicationID = mustGetEnv("APPLICATION_ID")
		SelfPinEnabled = strings.ToLower(os.Getenv("SELF_PIN_ENABLED")) == "true"
		HealthCheckAddr = os.Getenv("HEALTH_CHECK_ADDR")

		if s, ok := os.LookupEnv("EXCLUDED_CHANNELS"); ok && s != "" {
			ExcludedChannels = strings.Split(s, ",")
		}

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
		"SELF_PIN_ENABLED":  SelfPinEnabled,
		"EXCLUDED_CHANNELS": ExcludedChannels,
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

func IsExcludedChannel(id string) bool {
	for _, c := range ExcludedChannels {
		if c == id {
			return true
		}
	}

	return false
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
