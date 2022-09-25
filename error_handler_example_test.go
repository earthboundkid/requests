package requests_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/carlmjohnson/requests"
)

func ExampleValidationHandler() {
	logError := func(err error, req *http.Request, res *http.Response) {
		url := "<no url>"
		if req != nil {
			url = req.URL.String()
		}
		resCode := "---"
		if res != nil {
			resCode = res.Status
		}
		fmt.Printf("error kind=%s url=%s status=%s message=%v\n",
			requests.ErrorKindFrom(err), url, resCode, err)
	}
	var (
		body              string
		handledInvalidReq bool
		errBody           string
	)

	// All errors are sent to logErr.
	// If we fail validation because the response is a 404,
	// we send the body to errBody instead of body for separate
	// processing.
	err := requests.
		URL("http://example.com/404").
		ToString(&body).
		OnError(logError).
		OnValidationError(
			&handledInvalidReq, requests.ToString(&errBody)).
		Fetch(context.Background())
	if err != nil {
		fmt.Println(handledInvalidReq)
		fmt.Println(strings.Contains(errBody, "Example Domain"))
	}
	fmt.Println(strings.Contains(body, "Example Domain"))
	// Output:
	// error kind=ErrorKindValidator url=http://example.com/404 status=404 Not Found message=response error for http://example.com/404: unexpected status: 404
	// true
	// true
	// false
}
