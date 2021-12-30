package requests_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/carlmjohnson/requests"
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
	if err != nil {
		t.Fatal(err)
	}
	if h := headers.Headers["user-agent"]; h != "my-user/agent" {
		t.Fatalf("bad user agent: %q", h)
	}
}
