package requests_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/internal/be"
)

func TestErrorKind(t *testing.T) {
	var none requests.ErrorKind = -1
	kinds := []requests.ErrorKind{
		requests.ErrURL,
		requests.ErrRequest,
		requests.ErrTransport,
		requests.ErrValidator,
		requests.ErrHandler,
	}
	ctx := context.Background()
	res200 := requests.ReplayString("HTTP/1.1 200 OK\n\n")
	for _, tc := range []struct {
		ctx  context.Context
		want requests.ErrorKind
		b    *requests.Builder
	}{
		{ctx, none, requests.
			URL("").
			Transport(res200),
		},
		{ctx, requests.ErrURL, requests.
			URL("http://%2020").
			Transport(res200),
		},
		{ctx, requests.ErrURL, requests.
			URL("hello world").
			Transport(res200),
		},
		{ctx, none, requests.
			URL("http://world/#hello").
			Transport(res200),
		},
		{ctx, requests.ErrRequest, requests.
			URL("").
			Body(func() (io.ReadCloser, error) {
				return nil, errors.New("x")
			}).
			Transport(res200),
		},
		{ctx, requests.ErrRequest, requests.
			URL("").
			Method(" ").
			Transport(res200),
		},
		{nil, requests.ErrRequest, requests.
			URL("").
			Transport(res200),
		},
		{ctx, requests.ErrTransport, requests.
			URL("").
			Transport(requests.ReplayString("")),
		},
		{ctx, requests.ErrValidator, requests.
			URL("").
			Transport(requests.ReplayString("HTTP/1.1 404 Nope\n\n")),
		},
		{ctx, requests.ErrHandler, requests.
			URL("").
			Transport(res200).
			ToJSON(nil),
		},
	} {
		err := tc.b.Fetch(tc.ctx)
		for _, kind := range kinds {
			match := errors.Is(err, kind)
			be.Equal(t, kind == tc.want, match)
		}
		var askind = none
		be.Equal(t, tc.want != none, errors.As(err, &askind))
		be.Equal(t, tc.want, askind)
	}
}
