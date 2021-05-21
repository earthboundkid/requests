package requests_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/carlmjohnson/requests"
)

func Example() {
	// Simple GET into a buffer
	var buf bytes.Buffer
	err := requests.URL("http://example.com").
		Buffer(&buf).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to example.com:", err)
	}
	fmt.Println(strings.Contains(buf.String(), "Example Domain"))
	// Output:
	// true
}

func Example_bufio() {
	// read a response line by line for a sentinel
	found := false
	err := requests.URL("http://example.com").
		BufioReader(func(r *bufio.Reader) error {
			var err error
			for s := ""; err == nil; {
				if strings.Contains(s, "Example Domain") {
					found = true
					return nil
				}
				s, err = r.ReadString('\n')
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

type placeholder struct {
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

func Example_getJSON() {
	// GET a JSON object
	var post placeholder
	err := requests.URL("https://jsonplaceholder.typicode.com").
		Path("/posts/1").
		JSONUnmarshal(&post).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to jsonplaceholder.typicode.com:", err)
	}
	fmt.Println(post.Title)
	// Output:
	// sunt aut facere repellat provident occaecati excepturi optio reprehenderit
}

func Example_expectStatus() {
	// Expect a specific status code
	err := requests.URL("https://jsonplaceholder.typicode.com/posts/9001").
		AddValidator(requests.CheckStatus(404)).
		AddValidator(requests.MatchContentType("application/json")).
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
	err := requests.URL("/posts").
		Host("jsonplaceholder.typicode.com").
		JSONMarshal(&req).
		JSONUnmarshal(&res).
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
	err := requests.URL("https://postman-echo.com/get?hello=world").
		Param("param", "value").
		JSONUnmarshal(&params).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("problem with postman:", err)
	}
	fmt.Println(params.Args)
	// Output:
	// map[hello:world param:value]
}

func Example_header() {
	// Set a query parameter
	var params postman
	err := requests.URL("https://postman-echo.com/get?hello=world").
		UserAgent("bond/james-bond").
		ContentType("secret").
		Header("martini", "shaken").
		JSONUnmarshal(&params).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("problem with postman:", err)
	}
	fmt.Println(params.Headers["user-agent"])
	fmt.Println(params.Headers["content-type"])
	fmt.Println(params.Headers["martini"])
	// Output:
	// bond/james-bond
	// secret
	// shaken
}

func Example_rawBody() {
	// Post a raw body
	var data postman
	err := requests.URL("https://postman-echo.com/post").
		Bytes([]byte(`hello, world`), "text/plain").
		JSONUnmarshal(&data).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("problem with postman:", err)
	}
	fmt.Println(data.Data)
	// Output:
	// hello, world
}

func Example_formValue() {
	// Submit form values
	var echo postman
	err := requests.URL("https://postman-echo.com/put").
		Put().
		Form(url.Values{
			"hello": []string{"world"},
		}).
		JSONUnmarshal(&echo).
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
