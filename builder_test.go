package requests_test

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/internal/be"
)

func TestClone(t *testing.T) {
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
		be.Equal(t, "host.example3", req3.URL.Host)
		be.Equal(t, "/3", req3.URL.Path)
		be.Equal(t, "5", req3.Header.Get("b"))
		be.Equal(t, "6", req3.Header.Get("c"))
		be.Equal(t, "5", req3.URL.Query().Get("b"))
		be.Equal(t, "6", req3.URL.Query().Get("c"))
	})
}

func TestScheme(t *testing.T) {
	const res = `HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Date: Mon, 24 May 2021 18:48:50 GMT

An example response.`
	var s string
	const expected = `An example response.`
	var trans http.Transport
	trans.RegisterProtocol("string", requests.ReplayString(res))
	err := requests.
		URL("example").
		Scheme("string").
		Client(&http.Client{
			Transport: &trans,
		}).
		ToString(&s).
		Fetch(context.Background())
	be.NilErr(t, err)
	be.Equal(t, expected, s)
}

func TestPath(t *testing.T) {
	cases := map[string]struct {
		base   string
		paths  []string
		result string
	}{
		"base-only": {
			"example",
			[]string{},
			"https://example",
		},
		"base+abspath": {
			"https://example",
			[]string{"/a"},
			"https://example/a",
		},
		"multi-abs-paths": {
			"https://example",
			[]string{"/a", "/b/", "/c"},
			"https://example/c",
		},
		"base+rel-path": {
			"https://example/a/",
			[]string{"./b"},
			"https://example/a/b",
		},
		"base+rel-paths": {
			"https://example/a/",
			[]string{"./b/", "./c"},
			"https://example/a/b/c",
		},
		"rel-path": {
			"https://example/",
			[]string{"a/", "./b"},
			"https://example/a/b",
		},
		"base+multi-paths": {
			"https://example/a/",
			[]string{"b/", "c"},
			"https://example/a/b/c",
		},
		"base+slash+multi-paths": {
			"https://example/a/",
			[]string{"b/", "c"},
			"https://example/a/b/c",
		},
		"multi-root": {
			"https://example/",
			[]string{"a", "b", "c"},
			"https://example/c",
		},
		"dot-dot-paths": {
			"https://example/",
			[]string{"a/", "b/", "../c"},
			"https://example/a/c",
		},
		"more-dot-dot-paths": {
			"https://example/",
			[]string{"a/b/c/", "../d/", "../e"},
			"https://example/a/b/e",
		},
		"more-dot-dot-paths+rel-path": {
			"https://example/",
			[]string{"a/b/c/", "../d/", "../e/", "./f"},
			"https://example/a/b/e/f",
		},
		"even-more-dot-dot-paths+base": {
			"https://example/a/b/c/",
			[]string{"../../d"},
			"https://example/a/d",
		},
		"too-many-dot-dot-paths": {
			"https://example",
			[]string{"../a"},
			"https://example/a",
		},
		"too-many-dot-dot-paths+base": {
			"https://example/",
			[]string{"../a"},
			"https://example/a",
		},
		"last-abs-path-wins": {
			"https://example/a/",
			[]string{"b/", "c/", "/d"},
			"https://example/d",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			b := requests.URL(tc.base)
			for _, p := range tc.paths {
				b.Path(p)
			}
			r, err := b.Request(context.Background())
			be.NilErr(t, err)
			be.Equal(t, tc.result, r.URL.String())
		})
	}
}

func TestContentLength(t *testing.T) {
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
