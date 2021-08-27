package requests_test

import (
	"context"
	"fmt"
	"net/http"
	"testing/fstest"

	"github.com/carlmjohnson/requests"
)

func ExampleReplayFS() {
	fsys := fstest.MapFS{
		"fsys.example - MKIYDwjs.res.txt": &fstest.MapFile{
			Data: []byte(`HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Date: Mon, 24 May 2021 18:48:50 GMT

An example response.`),
		},
	}
	var s string
	const expected = `An example response.`
	if err := requests.
		URL("http://fsys.example").
		Client(&http.Client{
			Transport: requests.ReplayFS(fsys),
		}).
		ToString(&s).
		Fetch(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println(s == expected)
	// Output:
	// true
}
