package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/elliotwms/pinbot/internal/commands"
	"github.com/elliotwms/pinbot/internal/endpoint"
	"github.com/winebarrel/secretlamb"
)

const (
	envToken     = "PARAM_DISCORD_TOKEN"
	envPublicKey = "PARAM_DISCORD_PUBLIC_KEY"
)

var session *discordgo.Session
var publicKey []byte

func init() {
	params := secretlamb.MustNewParameters()

	if err := initDiscordSession(params); err != nil {
		panic(fmt.Errorf("failed to initialize Discord session: %v", err))
	}

	if err := initPublicKey(params); err != nil {
		panic(fmt.Errorf("failed to initialize public key: %v", err))
	}

	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}

func initDiscordSession(params *secretlamb.Parameters) error {
	p, err := params.GetWithDecryption(os.Getenv(envToken))
	if err != nil {
		return err
	}

	session, err = discordgo.New("Bot " + p.Parameter.Value)

	return err
}

func initPublicKey(params *secretlamb.Parameters) error {
	p, err := params.Get(os.Getenv(envPublicKey))
	if err != nil {
		return err
	}
	publicKey = []byte(p.Parameter.Value)

	return nil
}

func main() {
	h := endpoint.
		New(session).
		WithPublicKey(publicKey).
		WithApplicationCommand(commands.Pin.Name, commandhandlers.PinMessageCommandHandler)

	lambda.Start(h.Handle)
}
