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

const testGuildName = "Pinbot Integration Testing"

var (
	session     *discordgo.Session
	testGuildID string
)

var log = logrus.New()

func TestMain(m *testing.M) {
	fakediscord.Configure("http://localhost:8080/")

	_ = os.Setenv("TOKEN", "token")
	_ = os.Setenv("APPLICATION_ID", "appid")

	config.Configure()
	// enable testing with a single bot by allowing self-pins
	config.SelfPinEnabled = true

	// add additional testing permissions
	config.Permissions = config.DefaultPermissions |
		discordgo.PermissionManageChannels |
		discordgo.PermissionManageMessages

	openSession()

	code := m.Run()

	closeSession()

	os.Exit(code)
}

func openSession() {
	var err error
	session, err = discordgo.New(fmt.Sprintf("Bot %s", config.Token))
	if err != nil {
		panic(err)
	}

	if os.Getenv("TEST_DEBUG") != "" {
		session.LogLevel = discordgo.LogDebug
		session.Debug = true
	}

	session.Identify.Intents = config.Intents

	if err := session.Open(); err != nil {
		panic(err)
	}

	createGuild()
}

func createGuild() {
	guild, err := session.GuildCreate(testGuildName)
	if err != nil {
		panic(err)
	}

	testGuildID = guild.ID
}

func closeSession() {
	if err := session.Close(); err != nil {
		panic(err)
	}
}
