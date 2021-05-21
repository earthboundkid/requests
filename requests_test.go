package requests_test

import (
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
	s := buf.String()
	fmt.Println(strings.Contains(s, "Example Domain"))

	// GET a JSON object
	type placeholder struct {
		ID     int    `json:"id,omitempty"`
		Title  string `json:"title"`
		Body   string `json:"body"`
		UserID int    `json:"userId"`
	}
	var post placeholder
	err = requests.URL("https://jsonplaceholder.typicode.com").
		Path("/posts/1").
		JSONUnmarshal(&post).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to jsonplaceholder.typicode.com:", err)
	}
	fmt.Println(post.Title)

	// Expect a specific status code
	err = requests.URL("https://jsonplaceholder.typicode.com/posts/9001").
		Validate(requests.ChainHandlers(
			requests.CheckStatus(404),
			requests.MatchContentType("application/json"),
		)).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("should be a 404:", err)
	}

	// POST a JSON object and parse the response
	var res placeholder
	req := placeholder{
		Title:  "foo",
		Body:   "baz",
		UserID: 1,
	}
	err = requests.URL("/posts").
		Host("jsonplaceholder.typicode.com").
		JSONMarshal(&req).
		JSONUnmarshal(&res).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to jsonplaceholder.typicode.com:", err)
	}
	fmt.Println(res)

	// Examples with the Postman echo server
	type postman struct {
		Args map[string]string `json:"args"`
		Data string            `json:"data"`
		JSON map[string]string `json:"json"`
	}

	// Set a query parameter
	var params postman
	err = requests.URL("https://postman-echo.com/get?hello=world").
		Param("param", "value").
		JSONUnmarshal(&params).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("problem with postman:", err)
	}
	fmt.Println(params.Args)

	// Post a raw body
	var data postman
	err = requests.URL("https://postman-echo.com/post").
		Bytes([]byte(`hello, world`), "text/plain").
		JSONUnmarshal(&data).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("problem with postman:", err)
	}
	fmt.Println(data.Data)

	// Submit form values
	var echo postman
	err = requests.URL("https://postman-echo.com/put").
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
	// true
	// sunt aut facere repellat provident occaecati excepturi optio reprehenderit
	// {101 foo baz 1}
	// map[hello:world param:value]
	// hello, world
	// map[hello:world]
}
