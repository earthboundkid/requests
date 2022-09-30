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
	res200 := requests.ReplayString("HTTP/1.1 200 OK\n\n")
	for _, tc := range []struct {
		want requests.ErrorKind
		b    *requests.Builder
	}{
		{requests.KindNone, requests.
			URL("").
			Transport(res200),
		},
		{requests.KindURLErr, requests.
			URL("http://%2020").
			Transport(res200),
		},
		{requests.KindURLErr, requests.
			URL("hello world").
			Transport(res200),
		},
		{requests.KindBodyGet, requests.
			URL("").
			Body(func() (io.ReadCloser, error) {
				return nil, errors.New("x")
			}).
			Transport(res200),
		},
		{requests.KindBadMethod, requests.
			URL("").
			Method(" ").
			Transport(res200),
		},
		{requests.KindNilContext, requests.
			URL("").
			Transport(res200),
		},
		{requests.KindConnectErr, requests.
			URL("").
			Transport(requests.ReplayString("")),
		},
		{requests.KindInvalid, requests.
			URL("").
			Transport(requests.ReplayString("HTTP/1.1 404 Nope\n\n")),
		},
		{requests.KindHandlerErr, requests.
			URL("").
			Transport(res200).
			ToJSON(nil),
		},
	} {
		ctx := context.Background()
		if tc.want == requests.KindNilContext {
			ctx = nil
		}
		err := tc.b.Fetch(ctx)
		be.Equal(t, tc.want, requests.ErrorKindFrom(err))
	}

	be.Equal(t,
		requests.KindUnknown,
		requests.ErrorKindFrom(errors.New("")))
}
