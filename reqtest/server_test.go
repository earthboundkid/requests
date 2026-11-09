package reqtest_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gophertools/requests"
	"github.com/gophertools/requests/internal/be"
	"github.com/gophertools/requests/reqtest"
)

func TestServer(t *testing.T) {
	// Create an httptest.Server for your project's router
	mux := http.NewServeMux()
	mux.HandleFunc("/greeting", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	mux.HandleFunc("/salutation", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Howdy, planet!")
	})

	srv := httptest.NewTestServer(t, mux)

	// Now test that the handler has the expected return values
	{
		var s string
		err := requests.
			New(reqtest.Server(srv)).
			Path("/greeting").
			ToString(&s).
			Fetch(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		be.Equal(t, "Hello, world!", s)
	}
	{
		var s string
		err := requests.
			New(reqtest.Server(srv)).
			Path("/salutation").
			ToString(&s).
			Fetch(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		be.Equal(t, "Howdy, planet!", s)
	}
}
