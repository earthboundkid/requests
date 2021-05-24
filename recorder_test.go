package requests_test

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/carlmjohnson/requests"
)

func TestRecord(t *testing.T) {
	cl := *http.DefaultClient
	cl.Transport = requests.Record(cl.Transport, "testdata")
	var s1, s2 string
	err := requests.URL("http://example.com").
		Client(&cl).
		ToString(&s1).
		Fetch(context.Background())
	if err != nil {
		log.Fatal("unexpected error", err)
	}
	cl.Transport = requests.Replay("testdata")
	err = requests.URL("http://example.com").
		Client(&cl).
		ToString(&s2).
		Fetch(context.Background())
	if err != nil {
		log.Fatal("unexpected error", err)
	}
	if s1 != s2 {
		log.Fatalf("%q != %q", s1, s2)
	}
}
