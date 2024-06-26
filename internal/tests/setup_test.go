package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/fakediscord/pkg/fakediscord"
	"github.com/elliotwms/pinbot/internal/config"
	"github.com/sirupsen/logrus"
)

var botSession, userSession *discordgo.Session

var testGuildID string

var log = logrus.New()

func TestMain(m *testing.M) {
	fakediscord.Configure("http://localhost:8080/")

	_ = os.Setenv("TOKEN", "token")
	_ = os.Setenv("APPLICATION_ID", "appid")

	config.Configure()

	botSession, userSession = openSession("bot"), openSession("user")

	must(createGuild())

	code := m.Run()

	must(closeSession(botSession), closeSession(userSession))

	os.Exit(code)
}

func openSession(token string) *discordgo.Session {
	s, err := discordgo.New(fmt.Sprintf("Bot %s", token))
	must(err)

	if os.Getenv("TEST_DEBUG") != "" {
		s.LogLevel = discordgo.LogDebug
		s.Debug = true
	}

	s.Identify.Intents = config.Intents

	must(s.Open())

	return s
}

func closeSession(s *discordgo.Session) error {
	return s.Close()
}

func createGuild() error {
	guild, err := botSession.GuildCreate("Test")
	if err != nil {
		return err
	}

	testGuildID = guild.ID
	return nil
}

func must(errs ...error) {
	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}
