package tests

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
)

const testGuildName = "Pinbot Integration Testing"

var (
	botToken           string
	testGuild          string
	session            *discordgo.Session
	shouldCleanupGuild bool
)

func TestMain(m *testing.M) {
	configure()
	openSession()
	defer closeSession()

	m.Run()
	os.Exit(0)
}

func configure() {
	botToken = mustGetEnv("TOKEN")
}

func mustGetEnv(e string) string {
	s := os.Getenv(e)
	if s == "" {
		panic(fmt.Errorf("missing env '%s'", e))
	}

	return s
}

func openSession() {
	var err error
	session, err = discordgo.New(fmt.Sprintf("Bot %s", botToken))
	if err != nil {
		panic(err)
	}

	if err := session.Open(); err != nil {
		panic(err)
	}

	testGuild = os.Getenv("TEST_GUILD_ID")
	if testGuild == "" {
		log.Println("'TEST_GUILD_ID' not provided. Deleting stale guilds and creating new guild")
		deleteStaleGuilds()
		createGuild()
	} else {
		log.Printf("Using test guild ID '%s'", testGuild)
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

		if t, err := guild.JoinedAt.Parse(); err != nil || time.Since(t) > 10*time.Minute {
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

	testGuild = guild.ID
	shouldCleanupGuild = true
}

func closeSession() {
	if shouldCleanupGuild {
		_, err := session.GuildDelete(testGuild)
		if err != nil {
			panic(err)
		}
	}

	if err := session.Close(); err != nil {
		panic(err)
	}
}
