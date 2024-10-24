package endpoint

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elliotwms/pinbot/internal/build"
	"github.com/winebarrel/secretlamb"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/bwmarrin/discordgo"
)

const (
	envToken        = "PARAM_DISCORD_TOKEN"
	headerSignature = "X-Signature-Ed25519"
	headerTimestamp = "X-Signature-Timestamp"

	announcementURL = "https://discord.com/channels/1159611808722726912/1290727059261493358/1298783111265648693"
)

type Endpoint struct {
	s         *discordgo.Session
	handlers  map[string]CommandHandler
	publicKey ed25519.PublicKey
}

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error

func New(publicKey ed25519.PublicKey) *Endpoint {
	return &Endpoint{
		publicKey: publicKey,
		handlers:  map[string]CommandHandler{},
	}
}

func (r *Endpoint) WithSession(s *discordgo.Session) *Endpoint {
	r.s = s

	return r
}

func (r *Endpoint) WithApplicationCommand(name string, handler CommandHandler) *Endpoint {
	r.handlers[name] = handler
	return r
}

func (r *Endpoint) Handle(_ context.Context, event *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}

	bs := []byte(event.Body)

	slog.Info(
		"Received request",
		slog.String("user_agent", event.RequestContext.HTTP.UserAgent),
		slog.String("method", event.RequestContext.HTTP.Method),
		slog.String("version", build.Version),
	)

	if err := r.verify(event); err != nil {
		slog.Error("Failed to verify signature", "error", err)
		return &events.LambdaFunctionURLResponse{
			StatusCode: http.StatusUnauthorized,
		}, nil
	}

	var i *discordgo.InteractionCreate
	if err := json.Unmarshal(bs, &i); err != nil {
		return nil, err
	}

	response, err := r.handleInteraction(i)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return &events.LambdaFunctionURLResponse{StatusCode: http.StatusAccepted}, nil
	}

	bs, err = json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return &events.LambdaFunctionURLResponse{
		StatusCode: http.StatusOK,
		Body:       string(bs),
	}, nil
}

func (r *Endpoint) verify(event *events.LambdaFunctionURLRequest) error {
	if len(r.publicKey) == 0 {
		return nil
	}

	headers := make(http.Header, len(event.Headers))
	for k, v := range event.Headers {
		headers.Add(k, v)
	}

	signature := headers.Get(headerSignature)
	if signature == "" {
		return errors.New("missing header X-Signature-Ed25519")
	}
	ts := headers.Get(headerTimestamp)
	if ts == "" {
		return errors.New("missing header X-Signature-Timestamp")
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("invalid signature: %w", err)
	}

	verify := append([]byte(ts), []byte(event.Body)...)

	if !ed25519.Verify(r.publicKey, verify, sig) {
		return errors.New("invalid signature")
	}

	return nil
}

func (r *Endpoint) handleInteraction(i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	slog.Info("Handling interaction", "type", i.Type, "interaction_id", i.ID)

	switch i.Type {
	case discordgo.InteractionPing:
		return &discordgo.InteractionResponse{Type: discordgo.InteractionResponsePong}, nil
	case discordgo.InteractionApplicationCommand:
		// respond ASAP using the interaction's token
		is, _ := discordgo.New("Bot " + i.Token)
		if err := is.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			},
		}); err != nil {
			return nil, fmt.Errorf("initial respond: %w", err)
		}

		s, err := r.session()
		if err != nil {
			return nil, err
		}

		data := i.ApplicationCommandData()

		h, ok := r.handlers[data.Name]
		if !ok {
			return nil, cleanupStaleCommand(s, i, data)
		}

		if err = h(s, i, data); err != nil {
			return nil, fmt.Errorf("handle: %w", err)
		}

		return nil, nil
	default:
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Unexpected interaction"},
		}, nil
	}
}

func cleanupStaleCommand(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	log := slog.With("name", data.Name, "id", data.ID, "interaction_id", i.ID, "guild_id", i.GuildID)

	log.Info("Handling stale interaction")
	content := "This command is no longer supported. See the Pinbot Discord for more details: " + announcementURL
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	if err != nil {
		return err
	}

	slog.Debug("Removing stale command")

	// remove stale guild command
	return s.ApplicationCommandDelete(i.AppID, i.GuildID, data.ID)
}

// session returns the Bot session, initialising it if non-existent
func (r *Endpoint) session() (*discordgo.Session, error) {
	if r.s != nil {
		return r.s, nil
	}

	slog.Debug("Initiating new session")
	defer slog.Debug("Session initiated")

	var err error
	r.s, err = initDiscordSession()
	return r.s, err
}

// initDiscordSession initialises the Discord Session using the token stored in param store
func initDiscordSession() (*discordgo.Session, error) {
	paramName := os.Getenv(envToken)
	if paramName == "" {
		return nil, fmt.Errorf("missing required environment variable %q", envToken)
	}

	slog.Debug("Retrieving token")
	p, err := secretlamb.MustNewParameters().GetWithDecryption(paramName)
	if err != nil {
		return nil, err
	}

	slog.Debug("Retrieved token")

	if p == nil || p.Parameter.Value == "" {
		return nil, fmt.Errorf("parameter empty")
	}

	return discordgo.New("Bot " + p.Parameter.Value)
}
