package requests_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"
	"testing/fstest"

	"github.com/carlmjohnson/requests"
)

func TestRecord(t *testing.T) {
	cl := *http.DefaultClient
	dir := t.TempDir()
	cl.Transport = requests.Record(http.DefaultTransport, dir)
	var s1, s2 string
	err := requests.URL("http://example.com").
		Client(&cl).
		ToString(&s1).
		Fetch(context.Background())
	if err != nil {
		log.Fatalln("unexpected error:", err)
	}
	cl.Transport = requests.Replay(dir)
	err = requests.URL("http://example.com").
		Client(&cl).
		ToString(&s2).
		Fetch(context.Background())
	if err != nil {
		log.Fatalln("unexpected error:", err)
	}
	if s1 != s2 {
		log.Fatalf("%q != %q", s1, s2)
	}
}

func ExampleReplayFS() {
	fsys := fstest.MapFS{
		"MKIYDwjs.res.txt": &fstest.MapFile{
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

func TestScheme(t *testing.T) {
	fsys := fstest.MapFS{
		"AvbxD0vL.res.txt": &fstest.MapFile{
			Data: []byte(`HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Date: Mon, 24 May 2021 18:48:50 GMT

An example response.`),
		},
	}
	var s string
	const expected = `An example response.`
	var trans http.Transport
	trans.RegisterProtocol("fsys", requests.ReplayFS(fsys))
	if err := requests.
		URL("example").
		Scheme("fsys").
		Client(&http.Client{
			Transport: &trans,
		}).
		ToString(&s).
		Fetch(context.Background()); err != nil {
		t.Fatal(err)
	}
	if s != expected {
		t.Fatalf("%q != %q", s, expected)
	}
}
