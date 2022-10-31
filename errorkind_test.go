package requests_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/internal/be"
)

func kind(err error) requests.ErrorKind {
	var e requests.ErrorKindError
	if errors.As(err, &e) {
		return e.Kind()
	}
	return requests.ErrorKindUnknown
}

func TestErrorKind(t *testing.T) {
	ctx := context.Background()
	res200 := requests.ReplayString("HTTP/1.1 200 OK\n\n")
	for _, tc := range []struct {
		ctx  context.Context
		want requests.ErrorKind
		b    *requests.Builder
	}{
		{ctx, requests.ErrorKindUnknown, requests.
			URL("").
			Transport(res200),
		},
		{ctx, requests.ErrorKindURL, requests.
			URL("http://%2020").
			Transport(res200),
		},
		{ctx, requests.ErrorKindURL, requests.
			URL("hello world").
			Transport(res200),
		},
		{ctx, requests.ErrorKindBodyGet, requests.
			URL("").
			Body(func() (io.ReadCloser, error) {
				return nil, errors.New("x")
			}).
			Transport(res200),
		},
		{ctx, requests.ErrorKindMethod, requests.
			URL("").
			Method(" ").
			Transport(res200),
		},
		{nil, requests.ErrorKindContext, requests.
			URL("").
			Transport(res200),
		},
		{ctx, requests.ErrorKindConnect, requests.
			URL("").
			Transport(requests.ReplayString("")),
		},
		{ctx, requests.ErrorKindValidator, requests.
			URL("").
			Transport(requests.ReplayString("HTTP/1.1 404 Nope\n\n")),
		},
		{ctx, requests.ErrorKindHandler, requests.
			URL("").
			Transport(res200).
			ToJSON(nil),
		},
	} {
		err := tc.b.Fetch(tc.ctx)
		be.Equal(t, tc.want, kind(err))
	}

	be.Equal(t,
		requests.ErrorKindUnknown,
		kind(errors.New("")))
}
