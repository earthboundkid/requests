package requests_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/carlmjohnson/requests"
)

func TestUserAgentTransport(t *testing.T) {
	var headers postman
	err := requests.
		URL("https://postman-echo.com/get").
		Client(&http.Client{
			Transport: requests.UserAgentTransport(http.DefaultClient.Transport, "my-user/agent"),
		}).
		ToJSON(&headers).
		Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if h := headers.Headers["user-agent"]; h != "my-user/agent" {
		t.Fatalf("bad user agent: %q", h)
	}
}
