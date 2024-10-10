package main

import (
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/elliotwms/pinbot/internal/commands"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/router"
)

func main() {
	s, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		panic(err)
	}

	h := router.
		New(s).
		WithApplicationCommand(commands.Pin.Name, commandhandlers.PinMessageCommandHandler)

	lambda.Start(h.Handle)
}
