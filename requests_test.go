package requests_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/carlmjohnson/requests"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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

func ExampleBuilder_ToBufioReader() {
	// read a response line by line for a sentinel
	found := false
	err := requests.
		URL("http://example.com").
		ToBufioReader(func(r *bufio.Reader) error {
			var err error
			for s := ""; err == nil; {
				if strings.Contains(s, "Example Domain") {
					found = true
					return nil
				}
				// read one line from response
				s, err = r.ReadString('\n')
			}
			if err == io.EOF {
				return nil
			}
			return err
		}).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to example.com:", err)
	}
	fmt.Println(found)
	// Output:
	// true
}

func ExampleBuilder_ToHTML() {
	var doc html.Node
	err := requests.
		URL("http://example.com").
		ToHTML(&doc).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to example.com:", err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.DataAtom == atom.A {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					fmt.Println("link:", attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(&doc)
	// Output:
	// link: https://www.iana.org/domains/example
}

type placeholder struct {
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

func Example_getJSON() {
	// GET a JSON object
	var post placeholder
	err := requests.
		URL("https://jsonplaceholder.typicode.com").
		Path("/posts/1").
		ToJSON(&post).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to jsonplaceholder.typicode.com:", err)
	}
	fmt.Println(post.Title)
	// Output:
	// sunt aut facere repellat provident occaecati excepturi optio reprehenderit
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
		ContentType("secret").
		Header("martini", "shaken").
		ToJSON(&headers).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("problem with postman:", err)
	}
	fmt.Println(headers.Headers["user-agent"])
	fmt.Println(headers.Headers["content-type"])
	fmt.Println(headers.Headers["martini"])
	// Output:
	// bond/james-bond
	// secret
	// shaken
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

func ExampleHasStatusErr() {
	err := requests.
		URL("http://example.com/404").
		CheckStatus(200).
		Fetch(context.Background())
	if requests.HasStatusErr(err, 404) {
		fmt.Println("got a 404")
	}
	// Output:
	// got a 404
}

func ExampleBuilder_Peek() {
	// Check that a response has a doctype
	const doctype = "<!doctype html>"
	var s string
	err := requests.
		URL("http://example.com").
		Peek(len(doctype), func(b []byte) error {
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

func TestClone(t *testing.T) {
	{
		rb1 := requests.
			URL("example.com").
			Header("a", "1").
			Header("b", "2").
			Param("a", "1").
			Param("b", "2")
		rb2 := rb1.Clone().
			Host("host.example").
			Header("b", "3").
			Header("c", "4").
			Param("b", "3").
			Param("c", "4")
		req1, err := rb1.Request(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if req1.URL.Host != "example.com" {
			t.Fatalf("bad host: %v", req1.URL)
		}
		if req1.Header.Get("b") != "2" || req1.Header.Get("c") != "" {
			t.Fatalf("bad header: %v", req1.URL)
		}
		if q := req1.URL.Query(); q.Get("b") != "2" || q.Get("c") != "" {
			t.Fatalf("bad query: %v", req1.URL)
		}
		req2, err := rb2.Request(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if req2.URL.Host != "host.example" {
			t.Fatalf("bad host: %v", req2.URL)
		}
		if req2.Header.Get("b") != "3" || req2.Header.Get("c") != "4" {
			t.Fatalf("bad header: %v", req2.URL)
		}
		if q := req2.URL.Query(); q.Get("b") != "3" || q.Get("c") != "4" {
			t.Fatalf("bad query: %v", req2.URL)
		}
	}
	{
		rb1 := new(requests.Builder).
			Host("example.com").
			Header("a", "1").
			Header("b", "2").
			Param("a", "1").
			Param("b", "2")
		rb2 := rb1.Clone().
			Host("host.example").
			Header("b", "3").
			Header("c", "4").
			Param("b", "3").
			Param("c", "4")
		req1, err := rb1.Request(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if req1.URL.Host != "example.com" {
			t.Fatalf("bad host: %v", req1.URL)
		}
		if req1.Header.Get("b") != "2" || req1.Header.Get("c") != "" {
			t.Fatalf("bad header: %v", req1.URL)
		}
		if q := req1.URL.Query(); q.Get("b") != "2" || q.Get("c") != "" {
			t.Fatalf("bad query: %v", req1.URL)
		}
		req2, err := rb2.Request(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if req2.URL.Host != "host.example" {
			t.Fatalf("bad host: %v", req2.URL)
		}
		if req2.Header.Get("b") != "3" || req2.Header.Get("c") != "4" {
			t.Fatalf("bad header: %v", req2.URL)
		}
		if q := req2.URL.Query(); q.Get("b") != "3" || q.Get("c") != "4" {
			t.Fatalf("bad query: %v", req2.URL)
		}
	}
}
