package tests

import (
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	_, when, then := NewPingStage(t)

	when.a_ping_is_sent()

	then.the_status_code_should_be(http.StatusOK).and().
		a_pong_should_be_received()
}
