package requests_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/carlmjohnson/requests"
)

func ExampleNoFollow() {
	cl := *http.DefaultClient
	cl.CheckRedirect = requests.NoFollow

	var h http.Header
	if err := requests.
		URL("https://httpbingo.org/redirect/1").
		Client(&cl).
		CheckStatus(http.StatusFound).
		Handle(func(res *http.Response) error {
			h = res.Header
			return nil
		}).
		Fetch(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println(h.Get("Location"))
	// Output:
	// /get
}

func ExampleMaxFollow() {
	cl := *http.DefaultClient
	cl.CheckRedirect = requests.MaxFollow(1)

	var h http.Header
	if err := requests.
		URL("https://httpbingo.org/redirect/2").
		Client(&cl).
		CheckStatus(http.StatusFound).
		Handle(func(res *http.Response) error {
			h = res.Header
			return nil
		}).
		Fetch(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println(h.Get("Location"))
	// Output:
	// /get
}
