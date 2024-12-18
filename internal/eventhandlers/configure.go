package eventhandlers

import "github.com/sirupsen/logrus"

const logFieldHandler = "handler"

// List lists all event handlers to be registered when the bot is set up
func List(l *logrus.Entry) []interface{} {
	return []interface{}{
		Ready(l.WithField(logFieldHandler, "Ready")),
		InteractionCreate(l.WithField(logFieldHandler, "InteractionCreate")),
	}
}
