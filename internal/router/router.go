package router

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/bwmarrin/discordgo"
)

type Router struct {
	s        *discordgo.Session
	handlers map[string]CommandHandler
}

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) (string, error)

func New(s *discordgo.Session) *Router {
	return &Router{
		s:        s,
		handlers: map[string]CommandHandler{},
	}
}

func (r *Router) WithApplicationCommand(name string, handler CommandHandler) *Router {
	r.handlers[name] = handler
	return r
}

func (r *Router) Handle(_ context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}

	// todo validate headers

	var i *discordgo.InteractionCreate
	if err := json.Unmarshal([]byte(event.Body), &i); err != nil {
		return nil, err
	}

	response, err := r.handleInteraction(i)
	if err != nil {
		return nil, err
	}

	bs, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(bs),
	}, nil
}

func (r *Router) handleInteraction(i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
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
