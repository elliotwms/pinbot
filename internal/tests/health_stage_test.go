package tests

import (
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/config"
	"github.com/elliotwms/pinbot/internal/pinbot"
	"github.com/phayes/freeport"
	"github.com/sirupsen/logrus"
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
	log := logrus.New()

	p, err := freeport.GetFreePort()
	require.NoError(t, err)

	s := &HealthStage{
		t:       t,
		session: session,
		require: require.New(t),
		port:    strconv.Itoa(p),
	}

	done := make(chan os.Signal, 1)

	go func() {
		bot := pinbot.New(config.ApplicationID, session, log)
		s.require.NoError(bot.WithHealthCheck(":" + s.port).Run(done))
	}()

	t.Cleanup(func() {
		done <- os.Interrupt
	})

	return s, s, s
}

func (s *HealthStage) the_bot_is_running() *HealthStage {
	return s // no-op
}

func (s *HealthStage) a_health_check_request_is_sent() {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:"+s.port+"/v1/health", nil)
	s.require.NoError(err)
	s.res, err = http.DefaultClient.Do(req)
	s.require.NoError(err)
}

func (s *HealthStage) a_response_should_be_received_with_status_code(code int) {
	s.require.Equal(code, s.res.StatusCode)
}
