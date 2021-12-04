package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/storage"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	log.Info("Hello, World!")

	c, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	c.LogLevel = discordgo.LogDebug

	c.AddHandler(ready)
	c.AddHandler(react)
	c.AddHandler(guildCreate)
	c.AddHandler(channelCreate)
	c.AddHandler(channelUpdate)
	c.AddHandler(channelDelete)

	log.Info("Starting bot...")
	if err := c.Open(); err != nil {
		panic(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if err := c.Close(); err != nil {
		panic(err)
	}
}

func ready(s *discordgo.Session, _ *discordgo.Ready) {
	log.Info("I am ready for action")
	err := s.UpdateGameStatus(0, "Transcribing pins")
	if err != nil {
		log.WithError(err).Error("Could not update game status")
		return
	}
}

var guildChannels = sync.Map{}

func guildCreate(_ *discordgo.Session, e *discordgo.GuildCreate) {
	log.Info("Guild info received:", e.Name)

	gc, _ := guildChannels.LoadOrStore(e.Guild.ID, &storage.GuildChannels{})

	for _, c := range e.Channels {
		_, err := gc.(*storage.GuildChannels).Add(c)
		if err != nil {
			log.WithError(err).Error("Could not add channel")
			return
		}
	}
}

func channelCreate(_ *discordgo.Session, e *discordgo.ChannelCreate) {
	addChannel(e.Channel)
}

func channelUpdate(_ *discordgo.Session, e *discordgo.ChannelUpdate) {
	addChannel(e.Channel)
}

func channelDelete(_ *discordgo.Session, e *discordgo.ChannelDelete) {
	gc, _ := guildChannels.LoadOrStore(e.GuildID, &storage.GuildChannels{})

	_, err := gc.(*storage.GuildChannels).Add(e.Channel)
	if err != nil {
		log.WithError(err).Error("Could not add channel")
		return
	}
}

func addChannel(e *discordgo.Channel) {
	gc, _ := guildChannels.LoadOrStore(e.GuildID, &storage.GuildChannels{})

	_, err := gc.(*storage.GuildChannels).Add(e)
	if err != nil {
		log.WithError(err).Error("Could not add channel")
		return
	}
}

func react(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	log.WithField("emoji", e.Emoji.Name).Info("Received reaction")

	if e.Emoji.Name != "📌" {
		return
	}

	if testGuildID := os.Getenv("TEST_GUILD_ID"); testGuildID != "" && testGuildID != e.GuildID {
		log.Info("Skipping non-test channel")
		return
	}

	log.Info("PINNED")

	m, err := s.ChannelMessage(e.ChannelID, e.MessageID)
	if err != nil {
		log.WithError(err).Error("Could not get channel message")
		return
	}

	pinChannel, err := getTargetChannel(e.GuildID, e.ChannelID)
	if err != nil {
		log.WithError(err).Error("Could not get target channel")
		return
	}

	_, err = s.ChannelMessageSendEmbed(pinChannel, &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       e.Emoji.Name + " New Pin",
		Description: fmt.Sprintf("%s said: %s", m.Author.Mention(), m.Content),
		URL:         fmt.Sprintf("https://discord.com/channels/%s/%s/%s", e.GuildID, m.ChannelID, m.ID),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Pinned by %s", e.Member.User.String()),
		},
	})
	if err != nil {
		log.WithError(err).Error("Could not send message")
	}
}

// getTargetChannel returns the target pin channel for a given channel #channel in the following order:
// #channel-pins (a specific pin channel)
// #pins (a generic pin channel)
// #channel (the channel itself)
func getTargetChannel(guildID, channelID string) (string, error) {
	k, _ := guildChannels.LoadOrStore(guildID, &storage.GuildChannels{})
	gc, ok := k.(*storage.GuildChannels)
	if !ok {
		return "", errors.New("map did not contain type *storage.GuildChannels")
	}

	// get the channel
	var channel *discordgo.Channel

	for _, c := range gc.Channels {
		if c.ID == channelID {
			channel = c
			break
		}
	}

	if channel == nil {
		return "", errors.New("missing channel from map")
	}

	// check for #channel-pins
	for _, c := range gc.Channels {
		if c.Name == channel.Name+"-pins" {
			return c.ID, nil
		}
	}

	for _, c := range gc.Channels {
		if c.Name == "pins" {
			return c.ID, nil
		}
	}

	return channel.ID, nil
}
