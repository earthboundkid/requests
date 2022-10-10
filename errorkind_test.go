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
	ctx := context.Background()
	res200 := requests.ReplayString("HTTP/1.1 200 OK\n\n")
	for _, tc := range []struct {
		ctx  context.Context
		want requests.ErrorKind
		b    *requests.Builder
	}{
		{ctx, requests.KindNone, requests.
			URL("").
			Transport(res200),
		},
		{ctx, requests.KindURLErr, requests.
			URL("http://%2020").
			Transport(res200),
		},
		{ctx, requests.KindURLErr, requests.
			URL("hello world").
			Transport(res200),
		},
		{ctx, requests.KindBodyGetErr, requests.
			URL("").
			Body(func() (io.ReadCloser, error) {
				return nil, errors.New("x")
			}).
			Transport(res200),
		},
		{ctx, requests.KindMethodErr, requests.
			URL("").
			Method(" ").
			Transport(res200),
		},
		{nil, requests.KindContextErr, requests.
			URL("").
			Transport(res200),
		},
		{ctx, requests.KindConnectErr, requests.
			URL("").
			Transport(requests.ReplayString("")),
		},
		{ctx, requests.KindInvalidErr, requests.
			URL("").
			Transport(requests.ReplayString("HTTP/1.1 404 Nope\n\n")),
		},
		{ctx, requests.KindHandlerErr, requests.
			URL("").
			Transport(res200).
			ToJSON(nil),
		},
	} {
		err := tc.b.Fetch(tc.ctx)
		be.Equal(t, tc.want, requests.HasKindErr(err))
	}

	be.Equal(t,
		requests.KindUnknown,
		requests.HasKindErr(errors.New("")))
}
