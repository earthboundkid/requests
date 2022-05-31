package requests_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/internal/be"
)

func TestRecordReplay(t *testing.T) {
	dir := t.TempDir()

	var s1, s2 string
	err := requests.URL("http://example.com").
		Transport(requests.Record(http.DefaultTransport, dir)).
		ToString(&s1).
		Fetch(context.Background())
	be.NilErr(t, err)

	err = requests.URL("http://example.com").
		Transport(requests.Replay(dir)).
		ToString(&s2).
		Fetch(context.Background())
	be.NilErr(t, err)
	be.Equal(t, s1, s2)
}

func TestCaching(t *testing.T) {
	dir := t.TempDir()
	hasRun := false
	content := "some content"
	var onceTrans requests.RoundTripFunc = func(req *http.Request) (res *http.Response, err error) {
		be.False(t, hasRun)
		hasRun = true
		res = &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(content)),
		}
		return
	}
	trans := requests.Caching(onceTrans, dir)
	var s1, s2 string
	err := requests.URL("http://example.com").
		Transport(trans).
		ToString(&s1).
		Fetch(context.Background())
	be.NilErr(t, err)
	err = requests.URL("http://example.com").
		Transport(trans).
		ToString(&s2).
		Fetch(context.Background())
	be.NilErr(t, err)
	be.Equal(t, content, s1)
	be.Equal(t, s1, s2)

	entries, err := os.ReadDir(dir)
	be.NilErr(t, err)
	be.Equal(t, 2, len(entries))
}
