package reqtest

import (
	"net/http/httptest"

	"github.com/gophertools/requests"
)

// Server takes an httptest.Server and returns a requests.Config
// which sets the requests.Builder's BaseURL to s.URL
// and the requests.Builder's Client to s.Client().
func Server(s *httptest.Server) requests.Config {
	return func(rb *requests.Builder) {
		// In memory servers don't set s.URL until after s.Client is called,
		// so call that first.
		cl := s.Client()
		rb.
			BaseURL(s.URL).
			Client(cl)
	}
}
