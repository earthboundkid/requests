package reqtest

import (
	"net/http/httptest"

	"github.com/carlmjohnson/requests"
)

// Server takes an httptest.Server and returns a requests.Config
// which sets the requests.Builder's BaseURL to s.URL
// and the requests.Builder's Client to s.Client().
func Server(s *httptest.Server) requests.Config {
	return requests.TestServerConfig(s)
}
