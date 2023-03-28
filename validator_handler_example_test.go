package requests_test

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/carlmjohnson/requests"
)

func ExampleValidatorHandler() {
	var (
		regularBody string
		errBody     string
	)

	// If we fail validation because the response is a 404,
	// we handle the body with errBody instead of regularBody
	// for separate processing.
	err := requests.
		URL("http://example.com/404").
		ToString(&regularBody).
		AddValidator(
			requests.ValidatorHandler(
				requests.DefaultValidator,
				requests.ToString(&errBody),
			)).
		Fetch(context.Background())
	switch {
	case errors.Is(err, requests.ErrInvalidHandled):
		fmt.Println("got errBody:",
			strings.Contains(errBody, "Example Domain"))
	case err != nil:
		fmt.Println("unexpected error", err)
	case err == nil:
		fmt.Println("unexpected success")
	}

	fmt.Println("got regularBody:", strings.Contains(regularBody, "Example Domain"))
	// Output:
	// got errBody: true
	// got regularBody: false
}
