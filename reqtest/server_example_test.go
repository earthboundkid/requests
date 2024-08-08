package reqtest_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/reqtest"
)

func ExampleServer() {
	// Create an httptest.Server for your project's router
	mux := http.NewServeMux()
	mux.HandleFunc("/greeting", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	mux.HandleFunc("/salutation", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Howdy, planet!")
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	// Now test that the handler has the expected return values
	{
		var s string
		err := requests.
			New(reqtest.Server(srv)).
			Path("/greeting").
			ToString(&s).
			Fetch(context.Background())
		if err != nil {
			fmt.Println("Error!", err)
		}
		fmt.Println(s) // Hello, world!
	}
	{
		var s string
		err := requests.
			New(reqtest.Server(srv)).
			Path("/salutation").
			ToString(&s).
			Fetch(context.Background())
		if err != nil {
			fmt.Println("Error!", err)
		}
		fmt.Println(s) // Howdy, planet!
	}
	// Output:
	// Hello, world!
	// Howdy, planet!
}
