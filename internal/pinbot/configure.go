package pinbot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/eventhandlers"
	"github.com/sirupsen/logrus"
)

const logFieldHandler = "handler"

// RegisterHandlers registers the bot event handlers to an established session
func RegisterHandlers(session *discordgo.Session, l *logrus.Entry) func() {
	closers := []func(){
		session.AddHandler(eventhandlers.Ready(l.WithField(logFieldHandler, "Ready"))),
		session.AddHandler(eventhandlers.MessageReactionAdd(l.WithField(logFieldHandler, "MessageReactionAdd"))),
		session.AddHandler(eventhandlers.GuildCreate(l.WithField(logFieldHandler, "GuildCreate"))),
		session.AddHandler(eventhandlers.InteractionCreate(l.WithField(logFieldHandler, "InteractionCreate"))),
		session.AddHandler(eventhandlers.ChannelPinsUpdate(l.WithField(logFieldHandler, "ChannelPinsUpdate"))),
	}

	return func() {
		l.Debugf("Deregistering handlers (count: %d)", len(closers))
		for _, closer := range closers {
			closer()
		}
	}
}
