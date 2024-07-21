package reqtest

import (
	"net/http/httptest"

	"github.com/carlmjohnson/requests"
)

// TestServerConfig returns a Config
// which sets the Builder's BaseURL to s.URL
// and the Builder's Client to s.Client().
func TestServerConfig(s *httptest.Server) requests.Config {
	return requests.TestServerConfig(s)
}
