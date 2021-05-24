package requests

import "net/http"

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
