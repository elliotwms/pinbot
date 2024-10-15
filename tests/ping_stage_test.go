package tests

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/endpoint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PingStage struct {
	t           *testing.T
	session     *discordgo.Session
	require     *require.Assertions
	handler     func(context.Context, *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)
	res         *events.APIGatewayProxyResponse
	assert      *assert.Assertions
	privateKey  ed25519.PrivateKey
	omitHeaders bool
}

func NewPingStage(t *testing.T) (*PingStage, *PingStage, *PingStage) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	s := &PingStage{
		t:          t,
		assert:     assert.New(t),
		require:    require.New(t),
		session:    session,
		handler:    endpoint.New(session).WithPublicKey(publicKey).Handle,
		privateKey: privateKey,
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

	ts := strconv.FormatInt(time.Now().Unix(), 10)
	sign := ed25519.Sign(s.privateKey, append([]byte(ts), bs...))

	req := &events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Body:       string(bs),
	}

	if !s.omitHeaders {
		req.Headers = map[string]string{
			"X-Signature-Ed25519":   string(sign),
			"X-Signature-Timestamp": ts,
		}
	}

	s.res, err = s.handler(context.Background(), req)
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

func (s *PingStage) an_invalid_signature() {
	// trigger an invalid signature by changing the private key
	_, k, err := ed25519.GenerateKey(nil)
	s.require.NoError(err)

	s.privateKey = k
}

func (s *PingStage) request_will_omit_signature_headers() {
	s.omitHeaders = true
}
