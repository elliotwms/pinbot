package tests

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/elliotwms/fakediscord/pkg/fakediscord"
)

const testGuildName = "Pinbot Integration Testing"

var (
	session     *discordgo.Session
	testGuildID string
)

var node *snowflake.Node

func TestMain(m *testing.M) {
	fakediscord.Configure("http://localhost:8080/")

	if os.Getenv("TEST_DEBUG") != "" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	node, _ = snowflake.NewNode(0)

	openSession("bot")

	code := m.Run()

	closeSession()

	os.Exit(code)
}

func openSession(token string) {
	var err error
	session, err = discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		panic(err)
	}

	if os.Getenv("TEST_DEBUG") != "" {
		session.LogLevel = discordgo.LogDebug
		session.Debug = true
	}

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
