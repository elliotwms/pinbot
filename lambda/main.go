package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/elliotwms/pinbot/internal/commands"
	"github.com/elliotwms/pinbot/internal/endpoint"
)

func main() {
	s, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		panic(err)
	}

	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	h := endpoint.
		New(s).
		WithPublicKey([]byte(os.Getenv("DISCORD_PUBLIC_KEY"))).
		WithApplicationCommand(commands.Pin.Name, commandhandlers.PinMessageCommandHandler)

	lambda.Start(h.Handle)
}
