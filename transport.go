package requests

import (
	"bufio"
	"net/http"
	"strings"
)

// RoundTripFunc is an adaptor to use a function as an http.RoundTripper.
type RoundTripFunc func(req *http.Request) (res *http.Response, err error)

// RoundTrip implements http.RoundTripper.
func (rtf RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return rtf(r)
}

// ReplayString returns an http.RoundTripper that always responds with a
// request built from rawResponse. It is intended for use in one-off tests.
func ReplayString(rawResponse string) http.RoundTripper {
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		r := bufio.NewReader(strings.NewReader(rawResponse))
		res, err = http.ReadResponse(r, req)
		return
	})
}

// UserAgentTransport returns a wrapped http.RoundTripper that sets the User-Agent header on requests to s.
func UserAgentTransport(rt http.RoundTripper, s string) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		r2 := *req
		r2.Header = r2.Header.Clone()
		r2.Header.Set("User-Agent", s)
		return rt.RoundTrip(&r2)
	})
}
