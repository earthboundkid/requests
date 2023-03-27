package requests_test

import (
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/carlmjohnson/requests"
)

func ExampleNew() {
	// Suppose all requests in your project need some common options set.
	// First, define a Config function in your project...
	myProjectConfig := func(rb *requests.Builder) {
		rb.
			BaseURL("http://example.com").
			UserAgent("myproj/1.0").
			Accept("application/vnd.myproj+json;charset=utf-8")
	}

	// Then build your requests using that Config as the base Builder.
	var s string
	err := requests.
		New(myProjectConfig).
		Path("/").
		Param("some_param", "some-value").
		ToString(&s).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("my project fetch failed", err)
	}
	fmt.Println(strings.Contains(s, "Example Domain"))
	// Output:
	// true
}

func ExampleGzipConfig() {
	var echo postman
	err := requests.
		URL("https://postman-echo.com/post").
		ContentType("text/plain").
		Config(requests.GzipConfig(
			gzip.DefaultCompression,
			func(gw *gzip.Writer) error {
				_, err := gw.Write([]byte(`hello, world`))
				return err
			})).
		ToJSON(&echo).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("problem with postman:", err)
	}
	fmt.Println(echo.Data)
	// Output:
	// hello, world
}

func ExampleTestServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/greeting", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	mux.HandleFunc("/salutation", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Howdy, planet!")
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	{
		var s string
		err := requests.
			New(requests.TestServer(srv)).
			Path("/greeting").
			ToString(&s).
			Fetch(context.Background())
		if err != nil {
			fmt.Println("Error!", err)
		}
		fmt.Println(s)
	}
	{
		var s string
		err := requests.
			New(requests.TestServer(srv)).
			Path("/salutation").
			ToString(&s).
			Fetch(context.Background())
		if err != nil {
			fmt.Println("Error!", err)
		}
		fmt.Println(s)
	}
	// Output:
	// Hello, world!
	// Howdy, planet!
}
