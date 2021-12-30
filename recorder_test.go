package requests_test

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/carlmjohnson/requests"
)

func TestRecordReplay(t *testing.T) {
	dir := t.TempDir()

	var s1, s2 string
	err := requests.URL("http://example.com").
		Transport(requests.Record(http.DefaultTransport, dir)).
		ToString(&s1).
		Fetch(context.Background())
	if err != nil {
		log.Fatalln("unexpected error:", err)
	}

	err = requests.URL("http://example.com").
		Transport(requests.Replay(dir)).
		ToString(&s2).
		Fetch(context.Background())
	if err != nil {
		log.Fatalln("unexpected error:", err)
	}
	if s1 != s2 {
		log.Fatalf("%q != %q", s1, s2)
	}
}
