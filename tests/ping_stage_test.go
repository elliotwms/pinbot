package tests

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

type PingStage struct {
	t       *testing.T
	session *discordgo.Session
	require *require.Assertions
	handler func(context.Context, *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)
	res     *events.APIGatewayProxyResponse
	assert  *assert.Assertions
}

func NewPingStage(t *testing.T) (*PingStage, *PingStage, *PingStage) {
	s := &PingStage{
		t:       t,
		assert:  assert.New(t),
		require: require.New(t),
		session: session,
		handler: router.New(session).Handle,
	}

	return s, s, s
}

func (s *PingStage) and() *PingStage {
	return s
}

func (s *PingStage) a_ping_is_sent() *PingStage {
	bs, err := json.Marshal(&discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionPing,
		},
	})
	s.require.NoError(err)

	s.res, err = s.handler(context.Background(), &events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Body:       string(bs),
	})
	s.require.NoError(err)

	return s
}

func (s *PingStage) the_status_code_should_be(code int) *PingStage {
	s.assert.Equal(code, s.res.StatusCode)

	return s
}

func (s *PingStage) a_pong_should_be_received() {
	var res *discordgo.InteractionResponse

	err := json.Unmarshal([]byte(s.res.Body), &res)
	s.require.NoError(err)

	s.assert.Equal(discordgo.InteractionResponsePong, res.Type)
}
