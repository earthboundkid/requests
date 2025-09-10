package reqtest_test

import (
	"context"
	"testing"

	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/internal/be"
	"github.com/carlmjohnson/requests/reqtest"
)

func TestReplayJSON(t *testing.T) {
	err := requests.
		URL("http://example.com/api").
		Transport(reqtest.ReplayJSON(200, make(chan int))).
		Fetch(context.Background())
	be.Nonzero(t, err)
}
