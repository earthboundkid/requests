package requests_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/internal/be"
)

func TestClone(t *testing.T) {
	t.Parallel()
	t.Run("from URL", func(t *testing.T) {
		rb1 := requests.
			URL("http://example.com").
			Path("a/").
			Header("a", "1").
			Header("b", "2").
			Cookie("cookie", "base").
			Param("a", "1").
			Param("b", "2")
		rb2 := rb1.Clone().
			Host("host.example").
			Path("b").
			Header("b", "3").
			Header("c", "4").
			Cookie("cookie", "override").
			Param("b", "3").
			Param("c", "4")
		rb3 := rb1.Clone().
			Host("host.example3").
			Path("c").
			Header("b", "5").
			Header("c", "6").
			Cookie("alternate", "value").
			Param("b", "5").
			Param("c", "6")
		req1, err := rb1.Request(context.Background())
		be.NilErr(t, err)
		be.Equal(t, "example.com", req1.URL.Host)
		be.Equal(t, "/a/", req1.URL.Path)
		be.Equal(t, "2", req1.Header.Get("b"))
		be.Equal(t, "", req1.Header.Get("c"))
		be.Equal(t, "cookie=base", req1.Header.Get("Cookie"))
		be.Equal(t, "2", req1.URL.Query().Get("b"))
		be.Equal(t, "", req1.URL.Query().Get("c"))

		req2, err := rb2.Request(context.Background())
		be.NilErr(t, err)
		be.Equal(t, "host.example", req2.URL.Host)
		be.Equal(t, "/a/b", req2.URL.Path)
		be.Equal(t, "3", req2.Header.Get("b"))
		be.Equal(t, "4", req2.Header.Get("c"))
		be.Equal(t, "cookie=base; cookie=override", req2.Header.Get("Cookie"))
		be.Equal(t, "3", req2.URL.Query().Get("b"))
		be.Equal(t, "4", req2.URL.Query().Get("c"))

		req3, err := rb3.Request(context.Background())
		be.NilErr(t, err)
		be.Equal(t, "host.example3", req3.URL.Host)
		be.Equal(t, "/a/c", req3.URL.Path)
		be.Equal(t, "5", req3.Header.Get("b"))
		be.Equal(t, "6", req3.Header.Get("c"))
		be.Equal(t, "cookie=base; alternate=value", req3.Header.Get("Cookie"))
		be.Equal(t, "5", req3.URL.Query().Get("b"))
		be.Equal(t, "6", req3.URL.Query().Get("c"))
	})
	t.Run("from new", func(t *testing.T) {
		rb1 := new(requests.Builder).
			Host("example.com").
			Header("a", "1").
			Header("b", "2").
			Param("a", "1").
			Param("b", "2")
		rb2 := rb1.Clone().
			Host("host.example").
			Path("/2").
			Header("b", "3").
			Header("c", "4").
			Param("b", "3").
			Param("c", "4")
		rb3 := rb1.Clone().
			Host("host.example3").
			Path("/3").
			Header("b", "5").
			Header("c", "6").
			Param("b", "5").
			Param("c", "6")
		req1, err := rb1.Request(context.Background())
		be.NilErr(t, err)
		be.Equal(t, "example.com", req1.URL.Host)
		be.Equal(t, "", req1.URL.Path)
		be.Equal(t, "2", req1.Header.Get("b"))
		be.Equal(t, "", req1.Header.Get("c"))
		be.Equal(t, "2", req1.URL.Query().Get("b"))
		be.Equal(t, "", req1.URL.Query().Get("c"))

		req2, err := rb2.Request(context.Background())
		be.NilErr(t, err)
		be.Equal(t, "host.example", req2.URL.Host)
		be.Equal(t, "/2", req2.URL.Path)
		be.Equal(t, "3", req2.Header.Get("b"))
		be.Equal(t, "4", req2.Header.Get("c"))
		be.Equal(t, "3", req2.URL.Query().Get("b"))
		be.Equal(t, "4", req2.URL.Query().Get("c"))

		req3, err := rb3.Request(context.Background())
		be.NilErr(t, err)
		be.Equal(t, "host.example3", req3.URL.Host)
		be.Equal(t, "/3", req3.URL.Path)
		be.Equal(t, "5", req3.Header.Get("b"))
		be.Equal(t, "6", req3.Header.Get("c"))
		be.Equal(t, "5", req3.URL.Query().Get("b"))
		be.Equal(t, "6", req3.URL.Query().Get("c"))
	})
}

func TestScheme(t *testing.T) {
	u, err := requests.
		URL("example").
		Scheme("string").
		URL()
	be.NilErr(t, err)
	be.Equal(t, "string", u.Scheme)
	be.Equal(t, "example", u.Host)
	be.Equal(t, "string://example", u.String())
}

func TestPath(t *testing.T) {
	t.Parallel()
	for name, tc := range requests.PathCases {
		t.Run(name, func(t *testing.T) {
			var b requests.Builder
			b.BaseURL(tc.Base)
			for _, p := range tc.Paths {
				b.Path(p)
			}
			u, err := b.URL()
			be.NilErr(t, err)
			be.Equal(t, tc.Result, u.String())
		})
	}
}

func TestContentLength(t *testing.T) {
	t.Parallel()
	for _, n := range []int{0, 1, 10, 1000, 100_000} {
		req, err := requests.
			URL("http://example.com").
			BodyBytes(bytes.Repeat([]byte("a"), n)).
			Request(context.Background())
		be.NilErr(t, err)
		be.Equal(t, int64(n), req.ContentLength)
		req, err = requests.
			URL("http://example.com").
			BodyReader(
				strings.NewReader(strings.Repeat("a", n)),
			).
			Request(context.Background())
		be.NilErr(t, err)
		be.Equal(t, int64(n), req.ContentLength)
	}
	for _, obj := range []any{nil, 1, "x"} {
		req, err := requests.
			URL("http://example.com").
			BodyJSON(obj).
			Request(context.Background())
		be.NilErr(t, err)
		be.True(t, req.ContentLength > 0)
	}
	for _, qs := range []string{"", "a", "a=1"} {
		q, err := url.ParseQuery(qs)
		be.NilErr(t, err)
		req, err := requests.
			URL("http://example.com").
			BodyForm(q).
			Request(context.Background())
		be.NilErr(t, err)
		be.Equal(t, len(qs) > 0, req.ContentLength > 0)
	}
}

func TestRequestClosesBodyOnError(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		ctx    context.Context
		method string
		want   string
	}{
		{"invalid method", context.Background(), "bad method", "invalid method"},
		{"nil context", nil, "POST", "nil Context"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var file *os.File
			getter := requests.BodyFile("go.mod")
			req, err := requests.URL("https://example.com").
				Method(tc.method).
				Body(func() (io.ReadCloser, error) {
					body, err := getter()
					if err == nil {
						file = body.(*os.File)
						t.Cleanup(func() { file.Close() })
					}
					return body, err
				}).
				Request(tc.ctx)
			be.True(t, req == nil)
			be.True(t, errors.Is(err, requests.ErrRequest))
			be.True(t, strings.Contains(err.Error(), tc.want))
			if file == nil {
				t.Fatal("body file was not opened")
			}
			_, err = file.Read(make([]byte, 1))
			be.True(t, errors.Is(err, os.ErrClosed))
		})
	}
}

type requestBodyCloser struct {
	io.Reader
	closes int
	err    error
}

func (body *requestBodyCloser) Close() error {
	body.closes++
	return body.err
}

func TestRequestBodyOwnership(t *testing.T) {
	t.Parallel()
	closeErr := errors.New("close failed")
	for _, tc := range []struct {
		name    string
		method  string
		wantErr bool
	}{
		{"success", "POST", false},
		{"constructor error", "bad method", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := &requestBodyCloser{Reader: strings.NewReader("body"), err: closeErr}
			req, err := requests.URL("https://example.com").
				Method(tc.method).
				BodyReader(body).
				Request(context.Background())
			if tc.wantErr {
				be.True(t, req == nil)
				be.True(t, errors.Is(err, requests.ErrRequest))
				be.True(t, strings.Contains(err.Error(), "invalid method"))
				be.False(t, errors.Is(err, closeErr))
				be.Equal(t, 1, body.closes)
				return
			}
			be.NilErr(t, err)
			be.Equal(t, 0, body.closes)
			be.True(t, req.Body == body)
			be.Equal(t, closeErr, req.Body.Close())
			be.Equal(t, 1, body.closes)
		})
	}
}

func TestRequestErrorStopsBodyWriter(t *testing.T) {
	t.Parallel()
	done := make(chan error, 1)
	getter := requests.BodyWriter(func(w io.Writer) error {
		_, err := io.WriteString(w, "body")
		done <- err
		return err
	})
	req, err := requests.URL("https://example.com").
		Method("bad method").
		Body(func() (io.ReadCloser, error) {
			body, err := getter()
			if err == nil {
				t.Cleanup(func() { body.Close() })
			}
			return body, err
		}).
		Request(context.Background())
	be.True(t, req == nil)
	be.True(t, errors.Is(err, requests.ErrRequest))
	select {
	case err := <-done:
		be.True(t, errors.Is(err, io.ErrClosedPipe))
	case <-time.After(time.Second):
		t.Fatal("body writer did not stop after request construction failed")
	}
}
