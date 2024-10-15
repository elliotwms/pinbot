package endpoint

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/bwmarrin/discordgo"
)

type Endpoint struct {
	s         *discordgo.Session
	handlers  map[string]CommandHandler
	publicKey ed25519.PublicKey
}

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) (string, error)

func New(s *discordgo.Session) *Endpoint {
	return &Endpoint{
		s:        s,
		handlers: map[string]CommandHandler{},
	}
}

func (r *Endpoint) WithPublicKey(key ed25519.PublicKey) *Endpoint {
	r.publicKey = key

	return r
}

func (r *Endpoint) WithApplicationCommand(name string, handler CommandHandler) *Endpoint {
	r.handlers[name] = handler
	return r
}

func (r *Endpoint) Handle(_ context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}

	bs := []byte(event.Body)

	if !r.verify(event) {
		return &events.APIGatewayProxyResponse{
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

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(bs),
	}, nil
}

func (r *Endpoint) verify(event *events.APIGatewayProxyRequest) bool {
	if len(r.publicKey) == 0 {
		return true
	}

	sig, ok := event.Headers["X-Signature-Ed25519"]
	if !ok {
		return false
	}
	ts, ok := event.Headers["X-Signature-Timestamp"]
	if !ok {
		return false
	}
	verify := append([]byte(ts), []byte(event.Body)...)

	return ed25519.Verify(r.publicKey, verify, []byte(sig))
}

func (r *Endpoint) handleInteraction(i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	switch i.Type {
	case discordgo.InteractionPing:
		return &discordgo.InteractionResponse{Type: discordgo.InteractionResponsePong}, nil
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData()

		h, ok := r.handlers[data.Name]
		if !ok {
			return nil, fmt.Errorf("unknown command: %s", data.Name)
		}

		res, err := h(r.s, i, data)
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
