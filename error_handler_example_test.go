package requests_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/carlmjohnson/requests"
)

func ExampleBuilder_OnError() {
	logError := func(kind requests.ErrorKind, err error, req *http.Request, res *http.Response) error {
		url := "<no url>"
		if req != nil {
			url = req.URL.String()
		}
		resCode := "---"
		if res != nil {
			resCode = res.Status
		}
		fmt.Printf("[error] kind=%q url=%q status=%q message=%q\n",
			kind, url, resCode, err)
		return err
	}
	var (
		body    string
		errBody string
	)

	// All errors are sent to logErr.
	// If we fail validation because the response is a 404,
	// we send the body to errBody instead of body for separate
	// processing.
	err := requests.
		URL("http://example.com/404").
		ToString(&body).
		OnError(logError).
		OnValidatorError(
			requests.ToString(&errBody)).
		Fetch(context.Background())
	if err == requests.ErrInvalidHandled {
		fmt.Println("got errBody:",
			strings.Contains(errBody, "Example Domain"))
	} else if err != nil {
		fmt.Println("unknown error", err)
	}
	fmt.Println("got body:", strings.Contains(body, "Example Domain"))
	// Output:
	// [error] kind="ErrorKindValidator" url="http://example.com/404" status="404 Not Found" message="handled recovery from invalid response"
	// got errBody: true
	// got body: false
}
