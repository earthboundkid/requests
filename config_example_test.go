package requests_test

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/textproto"
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

func ExampleTestServerConfig() {
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
			New(requests.TestServerConfig(srv)).
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
			New(requests.TestServerConfig(srv)).
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

func ExampleBodyMultipart() {
	req, err := requests.
		URL("http://example.com").
		Config(requests.BodyMultipart("abc", func(multi *multipart.Writer) error {
			// CreateFormFile hardcodes the Content-Type as application/octet-stream
			w, err := multi.CreateFormFile("file", "en.txt")
			if err != nil {
				return err
			}
			_, err = io.WriteString(w, "Hello, World!")
			if err != nil {
				return err
			}
			// CreatePart is more flexible and lets you add headers
			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition", `form-data; name="file"; filename="jp.txt"`)
			h.Set("Content-Type", "text/plain; charset=utf-8")
			w, err = multi.CreatePart(h)
			if err != nil {
				panic(err)
			}
			_, err = io.WriteString(w, "こんにちは世界!")
			if err != nil {
				return err
			}
			return nil
		})).
		Request(context.Background())
	if err != nil {
		panic(err)
	}
	b, err := httputil.DumpRequest(req, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.ReplaceAll(string(b), "\r", ""))
	// Output:
	// POST / HTTP/1.1
	// Host: example.com
	// Content-Type: multipart/form-data; boundary=abc
	//
	// --abc
	// Content-Disposition: form-data; name="file"; filename="en.txt"
	// Content-Type: application/octet-stream
	//
	// Hello, World!
	// --abc
	// Content-Disposition: form-data; name="file"; filename="jp.txt"
	// Content-Type: text/plain; charset=utf-8
	//
	// こんにちは世界!
	// --abc--
}
