package tests

import (
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	_, when, then := NewPingStage(t)

	when.
		a_ping_is_sent()

	then.
		the_status_code_should_be(http.StatusOK).and().
		a_pong_should_be_received()
}

func TestPing_InvalidSignature(t *testing.T) {
	given, when, then := NewPingStage(t)

	given.
		an_invalid_signature()

	when.
		a_ping_is_sent()

	then.
		the_status_code_should_be(http.StatusUnauthorized)
}

func TestPing_MissingSignature(t *testing.T) {
	given, when, then := NewPingStage(t)

	given.
		request_will_omit_signature_headers()

	when.
		a_ping_is_sent()

	then.
		the_status_code_should_be(http.StatusUnauthorized)
}
