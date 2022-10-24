package requests

import (
	"net/http"
)

// ErrorHandler is a function accepted by Builder.OnError.
type ErrorHandler = func(error, *http.Request, *http.Response)

// ValidatorHandler converts a ResponseHandler into an ErrorHandler for invalid responses.
// The ResponseHandler only runs if the ErrorHandler encounters a validation error.
func ValidatorHandler(h ResponseHandler) ErrorHandler {
	return func(err error, req *http.Request, res *http.Response) {
		if res != nil && HasKindErr(err) == KindInvalidErr {
			h(res)
		}
	}
}
