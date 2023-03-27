package core

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/carlmjohnson/requests/internal/minitrue"
)

// NopCloser is like io.NopCloser(),
// but it is a concrete type so we can strip it out
// before setting a body on a request.
// See https://github.com/carlmjohnson/requests/discussions/49
type NopCloser struct {
	io.Reader
}

func RC(r io.Reader) NopCloser {
	return NopCloser{r}
}

func (NopCloser) Close() error { return nil }

var _ io.ReadCloser = NopCloser{}

// Request builds a new http.Request with its context set.
func (rb *RequestBuilder) Request(ctx context.Context, u *url.URL) (req *http.Request, err error) {
	var body io.Reader
	if rb.getBody != nil {
		if body, err = rb.getBody(); err != nil {
			return nil, err
		}
		if nopper, ok := body.(NopCloser); ok {
			body = nopper.Reader
		}
	}
	method := minitrue.First(rb.method,
		minitrue.Cond(rb.getBody != nil,
			http.MethodPost,
			http.MethodGet))

	req, err = http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.GetBody = rb.getBody

	for _, kv := range rb.headers {
		req.Header[http.CanonicalHeaderKey(kv.key)] = kv.values
	}
	for _, kv := range rb.cookies {
		req.AddCookie(&http.Cookie{
			Name:  kv.key,
			Value: kv.value,
		})
	}
	return req, nil
}
