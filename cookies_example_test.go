package requests_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/carlmjohnson/requests"
)

func ExampleNewCookieJar() {
	// Create a client that preserve cookies between requests
	myClient := *http.DefaultClient
	myClient.Jar = requests.NewCookieJar()
	// Use the client to make a request
	err := requests.
		URL("http://httpbin.org/cookies/set/chocolate/chip").
		Client(&myClient).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to httpbin.org:", err)
	}
	// Now check that cookies we got
	for _, cookie := range myClient.Jar.Cookies(&url.URL{
		Scheme: "http",
		Host:   "httpbin.org",
	}) {
		fmt.Println(cookie)
	}
	// And we'll see that they're reused on subsequent requests
	var cookies struct {
		Cookies map[string]string
	}
	err = requests.
		URL("http://httpbin.org/cookies").
		Client(&myClient).
		ToJSON(&cookies).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to httpbin.org:", err)
	}
	fmt.Println(cookies)
	// Output:
	// chocolate=chip
	// {map[chocolate:chip]}
}
