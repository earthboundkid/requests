package reqtest_test

import (
	"context"
	"fmt"
	"os"
	"testing/fstest"

	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/reqtest"
)

func ExampleReplayString() {
	const res = `HTTP/1.1 200 OK

An example response.`

	var s string
	const expected = `An example response.`
	if err := requests.
		URL("http://response.example").
		Transport(reqtest.ReplayString(res)).
		ToString(&s).
		Fetch(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println(s == expected)
	// Output:
	// true
}

func ExampleRecorder() {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}
	err = os.CopyFS(dir, fstest.MapFS{
		"fsys.example - MKIYDwjs.res.txt": &fstest.MapFile{
			Data: []byte(`HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Date: Mon, 24 May 2021 18:48:50 GMT

An example response.`),
		},
	})
	if err != nil {
		panic(err)
	}

	var s string
	const expected = `An example response.`
	if err := requests.
		URL("http://fsys.example").
		Transport(reqtest.Recorder(reqtest.ModeReplay, nil, dir)).
		ToString(&s).
		Fetch(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println(s == expected)
	// Output:
	// true
}
