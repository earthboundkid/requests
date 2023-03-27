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

func ExampleConfigure() {
	// Suppose all requests in your project need some common options set.
	// First, define a Config function in your project...
	myProjectConfig := func(rb *requests.Builder) {
		*rb = *requests.
			URL("http://example.com").
			UserAgent("myproj/1.0").
			Accept("application/vnd.myproj+json;charset=utf-8")
	}

	// Then build your requests using that Config as the base Builder.
	var s string
	err := requests.
		Configure(myProjectConfig).
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

func Example_GzipConfig() {
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
	h := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	}
	srv := httptest.NewServer(http.HandlerFunc(h))
	defer srv.Close()

	var s string
	err := requests.
		Configure(requests.TestServer(srv)).
		ToString(&s).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("Error!", err)
	}
	fmt.Println(s)
	// Output:
	// Hello, world!
}
