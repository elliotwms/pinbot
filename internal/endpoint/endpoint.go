package endpoint

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/winebarrel/secretlamb"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/bwmarrin/discordgo"
)

const (
	envToken        = "PARAM_DISCORD_TOKEN"
	headerUserAgent = "User-Agent"
	headerSignature = "X-Signature-Ed25519"
	headerTimestamp = "X-Signature-Timestamp"
)

type Endpoint struct {
	s             *discordgo.Session
	handlers      map[string]CommandHandler
	publicKey     ed25519.PublicKey
	applicationID string
}

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) (string, error)

func New(publicKey ed25519.PublicKey, applicationID string) *Endpoint {
	return &Endpoint{
		publicKey:     publicKey,
		applicationID: applicationID,
		handlers:      map[string]CommandHandler{},
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

	slog.Info("Received request", "user_agent", event.RequestContext.HTTP.UserAgent, "method", event.RequestContext.HTTP.Method)

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

	bs, err = json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return &events.LambdaFunctionURLResponse{
		StatusCode: 200,
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
	slog.Info("Handling interaction", "type", i.Type)

	switch i.Type {
	case discordgo.InteractionPing:
		return &discordgo.InteractionResponse{Type: discordgo.InteractionResponsePong}, nil
	case discordgo.InteractionApplicationCommand:
		s, err := r.session()
		if err != nil {
			return nil, err
		}

		data := i.ApplicationCommandData()

		h, ok := r.handlers[data.Name]
		if !ok {
			return nil, fmt.Errorf("unknown command: %s", data.Name)
		}

		res, err := h(s, i, data)
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: res},
		}, err
	default:
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Unexpected interaction"},
		}, nil
	}
}

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
