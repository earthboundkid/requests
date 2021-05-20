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
	return transport{rt, f}
}

type transport struct {
	rt http.RoundTripper
	f  func(*http.Request)
}

func (t transport) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	t.f(r)
	return t.rt.RoundTrip(r)
}
