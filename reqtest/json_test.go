package reqtest_test

import (
	"context"
	"testing"

	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/internal/be"
	"github.com/carlmjohnson/requests/reqtest"
)

func TestReplayJSON(t *testing.T) {
	// Marshals JSON
	var s string
	err := requests.
		URL("http://example.com/api").
		Transport(reqtest.ReplayJSON(200, "Hello")).
		ToString(&s).
		Fetch(context.Background())
	be.NilErr(t, err)
	be.Equal(t, `"Hello"`, s)

	// Returns marshal errors
	err = requests.
		URL("http://example.com/api").
		Transport(reqtest.ReplayJSON(200, make(chan int))).
		Fetch(context.Background())
	be.Nonzero(t, err)

	// Sets status code
	err = requests.
		URL("http://example.com/api").
		Transport(reqtest.ReplayJSON(500, "Error")).
		Fetch(context.Background())
	be.Nonzero(t, err)
	be.True(t, requests.HasStatusErr(err, 500))
}
