package requests

import (
	"bufio"
	"net/http"
	"strings"
)

// WrapTransport sets c.Transport to a WrapRoundTripper.
func WrapTransport(c *http.Client, f func(r *http.Request)) {
	c.Transport = WrapRoundTripper(c.Transport, f)
}

// WrapRoundTripper deep clones a request and passes it to f before calling the underlying
// RoundTripper. If rt is nil, it calls http.DefaultTransport.
func WrapRoundTripper(rt http.RoundTripper, f func(r *http.Request)) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return RoundTripFunc(func(r *http.Request) (*http.Response, error) {
		r = r.Clone(r.Context())
		f(r)
		return rt.RoundTrip(r)
	})
}

// RoundTripFunc implements http.RoundTripper.
type RoundTripFunc func(req *http.Request) (res *http.Response, err error)

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
