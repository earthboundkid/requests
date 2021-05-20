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
		Do(context.Background())
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
		Do(context.Background())
	if err != nil {
		fmt.Println("could not connect to jsonplaceholder.typicode.com:", err)
	}
	fmt.Println(post.Title)

	// Expect a specific status code
	err = requests.URL("https://jsonplaceholder.typicode.com/posts/9001").
		CheckStatus(404).
		Do(context.Background())
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
		Do(context.Background())
	if err != nil {
		fmt.Println("could not connect to jsonplaceholder.typicode.com:", err)
	}
	fmt.Println(res)

	// Submit form values
	echo := struct {
		JSON map[string]string `json:"json"`
	}{}
	err = requests.URL("https://postman-echo.com/put").
		Put().
		Form(url.Values{
			"hello": []string{"world"},
		}).
		JSONUnmarshal(&echo).
		Do(context.Background())
	if err != nil {
		fmt.Println("problem with postman", err)
	}
	fmt.Println(echo.JSON)
	// Output:
	// true
	// sunt aut facere repellat provident occaecati excepturi optio reprehenderit
	// {101 foo baz 1}
	// map[hello:world]
}
