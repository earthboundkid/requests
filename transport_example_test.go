package requests_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/carlmjohnson/requests"
)

func ExampleReplayString() {
	const res = `HTTP/1.1 200 OK

An example response.`

	var s string
	const expected = `An example response.`
	if err := requests.
		URL("http://response.example").
		Client(&http.Client{
			Transport: requests.ReplayString(res),
		}).
		ToString(&s).
		Fetch(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println(s == expected)
	// Output:
	// true
}
