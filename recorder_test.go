package requests_test

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/carlmjohnson/requests"
)

func TestRecordReplay(t *testing.T) {
	cl := *http.DefaultClient
	dir := t.TempDir()
	cl.Transport = requests.Record(http.DefaultTransport, dir)
	var s1, s2 string
	err := requests.URL("http://example.com").
		Client(&cl).
		ToString(&s1).
		Fetch(context.Background())
	if err != nil {
		log.Fatalln("unexpected error:", err)
	}
	cl.Transport = requests.Replay(dir)
	err = requests.URL("http://example.com").
		Client(&cl).
		ToString(&s2).
		Fetch(context.Background())
	if err != nil {
		log.Fatalln("unexpected error:", err)
	}
	if s1 != s2 {
		log.Fatalf("%q != %q", s1, s2)
	}
}
