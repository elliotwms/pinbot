package tests

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/fakediscord/pkg/fakediscord"
	"github.com/elliotwms/pinbot/internal/config"
)

const testGuildName = "Pinbot Integration Testing"

var (
	session            *discordgo.Session
	shouldCleanupGuild bool
)

func TestMain(m *testing.M) {
	if v := os.Getenv("FAKEDISCORD"); v != "" {
		fakediscord.Configure("http://localhost:8080/")

		_ = os.Setenv("TOKEN", "token")
		_ = os.Setenv("APPLICATION_ID", "appid")
	}

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

	session.LogLevel = discordgo.LogDebug

	session.Identify.Intents = config.Intents

	if err := session.Open(); err != nil {
		panic(err)
	}

	if config.TestGuildID == "" {
		log.Println("'TEST_GUILD_ID' not provided. Deleting stale guilds and creating new guild")
		deleteStaleGuilds()
		createGuild()
	} else {
		log.Printf("Using test guild ID '%s'", config.TestGuildID)
		log.Printf("Ensure the bot is installed in the test guild with the additional test permissions: %s", config.BuildInstallURL())
	}
}

func deleteStaleGuilds() {
	guilds, err := session.UserGuilds(0, "", "")
	if err != nil {
		panic(err)
	}

	for _, ug := range guilds {
		if ug.Name != testGuildName {
			continue
		}

		guild, err := session.Guild(ug.ID)
		if err != nil {
			panic(err)
		}

		if time.Since(guild.JoinedAt) > 10*time.Minute {
			log.Printf("Deleting stale guild '%s'", guild.ID)
			if _, err := session.GuildDelete(guild.ID); err != nil {
				return
			}
		}
	}
}

func createGuild() {
	guild, err := session.GuildCreate(testGuildName)
	if err != nil {
		panic(err)
	}

	config.TestGuildID = guild.ID
	shouldCleanupGuild = true
}

func closeSession() {
	if shouldCleanupGuild {
		_, _ = session.GuildDelete(config.TestGuildID)
	}

	if err := session.Close(); err != nil {
		panic(err)
	}
}
