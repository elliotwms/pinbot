package tests

import (
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {
	given, when, then := NewHealthStage(t)

	given.the_bot_is_running()
	when.a_health_check_request_is_sent()
	then.a_response_should_be_received_with_status_code(http.StatusOK)
}
