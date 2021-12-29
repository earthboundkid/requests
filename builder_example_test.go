package requests_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/carlmjohnson/requests"
)

func init() {
	http.DefaultClient.Transport = requests.Replay("testdata")
}

func Example() {
	// Simple GET into a string
	var s string
	err := requests.
		URL("http://example.com").
		ToString(&s).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to example.com:", err)
	}
	fmt.Println(strings.Contains(s, "Example Domain"))
	// Output:
	// true
}

func Example_getJSON() {
	// GET a JSON object
	id := 1
	var post placeholder
	err := requests.
		URL("https://jsonplaceholder.typicode.com").
		Pathf("/posts/%d", id).
		ToJSON(&post).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to jsonplaceholder.typicode.com:", err)
	}
	fmt.Println(post.Title)
	// Output:
	// sunt aut facere repellat provident occaecati excepturi optio reprehenderit
}

func Example_postJSON() {
	// POST a JSON object and parse the response
	var res placeholder
	req := placeholder{
		Title:  "foo",
		Body:   "baz",
		UserID: 1,
	}
	err := requests.
		URL("/posts").
		Host("jsonplaceholder.typicode.com").
		BodyJSON(&req).
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to jsonplaceholder.typicode.com:", err)
	}
	fmt.Println(res)
	// Output:
	// {101 foo baz 1}
}

func ExampleBuilder_ToBytesBuffer() {
	// Simple GET into a buffer
	var buf bytes.Buffer
	err := requests.
		URL("http://example.com").
		ToBytesBuffer(&buf).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to example.com:", err)
	}
	fmt.Println(strings.Contains(buf.String(), "Example Domain"))
	// Output:
	// true
}

func ExampleBuilder_ToWriter() {
	f, err := os.CreateTemp("", "*.example.html")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name()) // clean up

	err = requests.
		URL("http://example.com").
		ToWriter(f).
		Fetch(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if err = f.Close(); err != nil {
		log.Fatal(err)
	}
	stat, err := os.Stat(f.Name())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("file is %d bytes\n", stat.Size())

	// Output:
	// file is 1256 bytes
}

func ExampleBuilder_ToFile() {
	d, err := os.MkdirTemp("", "to_file_*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(d) // clean up

	exampleFilename := filepath.Join(d, "parent_dir", "example.txt")

	err = requests.
		URL("http://example.com").
		ToFile(exampleFilename).
		Fetch(context.Background())

	if err != nil {
		log.Fatal(err)
	}
	stat, err := os.Stat(exampleFilename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("file is %d bytes\n", stat.Size())

	// Output:
	// file is 1256 bytes
}

type placeholder struct {
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

func ExampleBuilder_Path() {
	// Add an ID to a base path
	id := 1
	var post placeholder
	err := requests.
		URL("https://jsonplaceholder.typicode.com/posts/").
		// inherits path /posts from baseurl
		Pathf("%d", id).
		// URL is now https://jsonplaceholder.typicode.com/posts/1
		ToJSON(&post).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to jsonplaceholder.typicode.com:", err)
	}
	fmt.Println(post.ID)
	// Output:
	// 1
}
func ExampleBuilder_CheckStatus() {
	// Expect a specific status code
	err := requests.
		URL("https://jsonplaceholder.typicode.com").
		Pathf("/posts/%d", 9001).
		CheckStatus(404).
		CheckContentType("application/json").
		Fetch(context.Background())
	if err != nil {
		fmt.Println("should be a 404:", err)
	} else {
		fmt.Println("OK")
	}
	// Output:
	// OK
}

func ExampleBuilder_CheckContentType() {
	// Expect a specific status code
	err := requests.
		URL("https://jsonplaceholder.typicode.com").
		Pathf("/posts/%d", 1).
		CheckContentType("application/bison").
		Fetch(context.Background())
	if err != nil {
		if re := new(requests.ResponseError); errors.As(err, &re) {
			fmt.Println("content-type was", re.Header.Get("Content-Type"))
		}
	}
	// Output:
	// content-type was application/json; charset=utf-8
}

// Examples with the Postman echo server
type postman struct {
	Args    map[string]string `json:"args"`
	Data    string            `json:"data"`
	Headers map[string]string `json:"headers"`
	JSON    map[string]string `json:"json"`
}

func Example_queryParam() {
	// Set a query parameter
	var params postman
	err := requests.
		URL("https://postman-echo.com/get?a=1&b=2").
		Param("b", "3").
		Param("c", "4").
		ToJSON(&params).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("problem with postman:", err)
	}
	fmt.Println(params.Args)
	// Output:
	// map[a:1 b:3 c:4]
}

func ExampleBuilder_Header() {
	// Set headers
	var headers postman
	err := requests.
		URL("https://postman-echo.com/get").
		UserAgent("bond/james-bond").
		BasicAuth("bondj", "007!").
		ContentType("secret").
		Header("martini", "shaken").
		ToJSON(&headers).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("problem with postman:", err)
	}
	fmt.Println(headers.Headers["user-agent"])
	fmt.Println(headers.Headers["authorization"])
	fmt.Println(headers.Headers["content-type"])
	fmt.Println(headers.Headers["martini"])
	// Output:
	// bond/james-bond
	// Basic Ym9uZGo6MDA3IQ==
	// secret
	// shaken
}

func ExampleBuilder_Bearer() {
	// We get a 401 response if no bearer token is provided
	err := requests.
		URL("http://httpbin.org/bearer").
		CheckStatus(http.StatusUnauthorized).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("problem with httpbin:", err)
	}
	// But our response is accepted when we provide a bearer token
	var res struct {
		Authenticated bool
		Token         string
	}
	err = requests.
		URL("http://httpbin.org/bearer").
		Bearer("whatever").
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("problem with httpbin:", err)
	}
	fmt.Println(res.Authenticated)
	fmt.Println(res.Token)
	// Output:
	// true
	// whatever
}

func ExampleBuilder_BodyBytes() {
	// Post a raw body
	var data postman
	err := requests.
		URL("https://postman-echo.com/post").
		BodyBytes([]byte(`hello, world`)).
		ContentType("text/plain").
		ToJSON(&data).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("problem with postman:", err)
	}
	fmt.Println(data.Data)
	// Output:
	// hello, world
}

func ExampleBuilder_BodyReader() {
	// temp file creation boilerplate
	dir, err := os.MkdirTemp("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up

	file := filepath.Join(dir, "tmpfile")
	if err := os.WriteFile(file, []byte(`hello, world`), 0666); err != nil {
		log.Fatal(err)
	}

	// suppose there is some file you want to stream from
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// send the raw file to server
	var echo postman
	err = requests.
		URL("https://postman-echo.com/post").
		ContentType("text/plain").
		BodyReader(f).
		ToJSON(&echo).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("problem with postman:", err)
	}
	fmt.Println(echo.Data)
	// Output:
	// hello, world
}

func ExampleBuilder_Config() {
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

func ExampleBuilder_BodyWriter() {
	var echo postman
	err := requests.
		URL("https://postman-echo.com/post").
		ContentType("text/plain").
		BodyWriter(func(w io.Writer) error {
			cw := csv.NewWriter(w)
			cw.Write([]string{"col1", "col2"})
			cw.Write([]string{"val1", "val2"})
			cw.Flush()
			return cw.Error()
		}).
		ToJSON(&echo).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("problem with postman:", err)
	}
	fmt.Printf("%q\n", echo.Data)
	// Output:
	// "col1,col2\nval1,val2\n"
}

func ExampleBuilder_BodyForm() {
	// Submit form values
	var echo postman
	err := requests.
		URL("https://postman-echo.com/put").
		Put().
		BodyForm(url.Values{
			"hello": []string{"world"},
		}).
		ToJSON(&echo).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("problem with postman:", err)
	}
	fmt.Println(echo.JSON)
	// Output:
	// map[hello:world]
}

func ExampleBuilder_BodyFile() {
	// Make a file to read from
	d, err := os.MkdirTemp("", "body_file_*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(d) // clean up

	exampleFilename := filepath.Join(d, "example.txt")
	exampleContent := `hello, world`
	if err = os.WriteFile(exampleFilename, []byte(exampleContent), 0644); err != nil {
		log.Fatal(err)
	}

	// Post a raw file
	var data postman
	err = requests.
		URL("https://postman-echo.com/post").
		BodyFile(exampleFilename).
		ContentType("text/plain").
		ToJSON(&data).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("problem with postman:", err)
	}
	fmt.Println(data.Data)
	// Output:
	// hello, world
}

func ExampleBuilder_CheckPeek() {
	// Check that a response has a doctype
	const doctype = "<!doctype html>"
	var s string
	err := requests.
		URL("http://example.com").
		CheckPeek(len(doctype), func(b []byte) error {
			if string(b) != doctype {
				return fmt.Errorf("missing doctype: %q", b)
			}
			return nil
		}).
		ToString(&s).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to example.com:", err)
	}
	fmt.Println(
		// Final result still has the prefix
		strings.HasPrefix(s, doctype),
		// And the full body
		strings.HasSuffix(s, "</html>\n"),
	)
	// Output:
	// true true
}

func ExampleBuilder_Transport() {
	const text = "Hello, from transport!"
	var myCustomTransport requests.RoundTripFunc = func(req *http.Request) (res *http.Response, err error) {
		res = &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(text)),
		}
		return
	}
	var s string
	err := requests.
		URL("x://transport.example").
		Transport(myCustomTransport).
		ToString(&s).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("transport failed:", err)
	}
	fmt.Println(s == text) // true
	// Output:
	// true
}
