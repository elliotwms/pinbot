package tests

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/bot"
	"github.com/elliotwms/pinbot/internal/config"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/require"
)

type HealthStage struct {
	t       *testing.T
	session *discordgo.Session
	require *require.Assertions
	res     *http.Response
	port    string
}

func NewHealthStage(t *testing.T) (*HealthStage, *HealthStage, *HealthStage) {
	p, err := freeport.GetFreePort()
	require.NoError(t, err)

	s := &HealthStage{
		t:       t,
		session: session,
		require: require.New(t),
		port:    strconv.Itoa(p),
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		b := bot.New(config.ApplicationID, session, log)
		s.require.NoError(b.WithHealthCheck(":" + s.port).Run(ctx))
	}()

	t.Cleanup(cancel)

	return s, s, s
}

func (s *HealthStage) the_bot_is_running() *HealthStage {
	return s // no-op
}

func (s *HealthStage) a_health_check_request_is_sent() {
	s.require.Eventually(func() bool {
		var err error
		s.res, err = http.Get("http://localhost:" + s.port + "/v1/health")

		return err == nil
	}, 1*time.Second, 10*time.Millisecond)
}

func (s *HealthStage) a_response_should_be_received_with_status_code(code int) {
	s.require.Equal(code, s.res.StatusCode)
}
