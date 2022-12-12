package requests_test

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/carlmjohnson/requests"
)

func ExampleBuilder_OnError() {
	logError := func(ep *requests.OnErrorParams) {
		fmt.Printf("[error] kind=%q method=%s url=%q status=%3d message=%q\n",
			ep.Kind(), ep.Method(), ep.URL(), ep.StatusCode(), ep.Error)
		ep.Error = fmt.Errorf("my app error: %w", ep.Error)
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
	if errors.Is(err, requests.ErrInvalidHandled) {
		fmt.Println("is a my-app-error:",
			strings.Contains(err.Error(), "my app error:"))
		fmt.Println("got errBody:",
			strings.Contains(errBody, "Example Domain"))
	} else if err != nil {
		fmt.Println("unknown error", err)
	}
	fmt.Println("got body:", strings.Contains(body, "Example Domain"))
	// Output:
	// [error] kind="ErrValidator" method=GET url="http://example.com/404" status=404 message="handled recovery from invalid response"
	// is a my-app-error: true
	// got errBody: true
	// got body: false
}
