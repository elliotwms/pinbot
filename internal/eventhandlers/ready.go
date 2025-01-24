package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/build"
	"log/slog"
)

func Ready(s *discordgo.Session, _ *discordgo.Ready) {
	slog.Info("I am ready for action")
	err := s.UpdateGameStatus(0, build.Version)
	if err != nil {
		slog.Error("Could not update game status", "error", err)
		return
	}
}
