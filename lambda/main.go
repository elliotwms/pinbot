package main

import (
	"encoding/hex"
	"log/slog"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/elliotwms/pinbot/internal/commands"
	"github.com/elliotwms/pinbot/internal/endpoint"
)

func init() {
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}

func main() {
	k, err := hex.DecodeString(os.Getenv("DISCORD_BOT_PUBLIC_KEY"))
	if err != nil {
		panic(err)
	}

	h := endpoint.
		New(k).
		WithApplicationCommand(commands.Pin.Name, commandhandlers.PinMessageCommandHandler)

	lambda.Start(h.Handle)
}
