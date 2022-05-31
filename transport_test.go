package requests_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/internal/be"
)

func TestUserAgentTransport(t *testing.T) {
	// Wrap an existing transport or use nil for http.DefaultTransport
	baseTrans := http.DefaultClient.Transport
	trans := requests.UserAgentTransport(baseTrans, "my-user/agent")

	var headers postman
	err := requests.
		URL("https://postman-echo.com/get").
		Transport(trans).
		ToJSON(&headers).
		Fetch(context.Background())
	be.NilErr(t, err)
	be.Equal(t, "my-user/agent", headers.Headers["user-agent"])
}
