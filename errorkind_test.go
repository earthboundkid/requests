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
	kind := requests.ErrUnknown
	setKind := func(ep *requests.OnErrorParams) {
		kind = ep.Kind()
	}
	ctx := context.Background()
	res200 := requests.ReplayString("HTTP/1.1 200 OK\n\n")
	for _, tc := range []struct {
		ctx  context.Context
		want requests.ErrorKind
		b    *requests.Builder
	}{
		{ctx, requests.ErrUnknown, requests.
			URL("").
			Transport(res200).
			OnError(setKind),
		},
		{ctx, requests.ErrURL, requests.
			URL("http://%2020").
			Transport(res200).
			OnError(setKind),
		},
		{ctx, requests.ErrURL, requests.
			URL("hello world").
			Transport(res200).
			OnError(setKind),
		},
		{ctx, requests.ErrUnknown, requests.
			URL("http://world/#hello").
			Transport(res200).
			OnError(setKind),
		},
		{ctx, requests.ErrRequest, requests.
			URL("").
			Body(func() (io.ReadCloser, error) {
				return nil, errors.New("x")
			}).
			Transport(res200).
			OnError(setKind),
		},
		{ctx, requests.ErrRequest, requests.
			URL("").
			Method(" ").
			Transport(res200).
			OnError(setKind),
		},
		{nil, requests.ErrRequest, requests.
			URL("").
			Transport(res200).
			OnError(setKind),
		},
		{ctx, requests.ErrConnect, requests.
			URL("").
			Transport(requests.ReplayString("")).
			OnError(setKind),
		},
		{ctx, requests.ErrValidator, requests.
			URL("").
			Transport(requests.ReplayString("HTTP/1.1 404 Nope\n\n")).
			OnError(setKind),
		},
		{ctx, requests.ErrHandler, requests.
			URL("").
			Transport(res200).
			ToJSON(nil).
			OnError(setKind),
		},
	} {
		err := tc.b.Fetch(tc.ctx)
		be.Equal(t, tc.want, kind)
		be.Equal(t, tc.want != requests.ErrUnknown, errors.Is(err, tc.want))
		var askind requests.ErrorKind
		be.Equal(t, tc.want != requests.ErrUnknown, errors.As(err, &askind))
		be.Equal(t, tc.want, askind)
		kind = requests.ErrUnknown
	}
}
