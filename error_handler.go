package requests

import (
	"net/http"
)

// ErrorHandler is a function accepted by Builder.OnError.
type ErrorHandler = func(error, *http.Request, *http.Response)

// ValidationHandler converts a ResponseHandler into an ErrorHandler for invalid responses. If the error handling succeeds, it sets ok to true.
// The handler only runs if the ErrorHandler encounters a validation error.
// If ok is nil, the ErrorHandler ignores it.
func ValidationHandler(ok *bool, h ResponseHandler) ErrorHandler {
	if ok == nil {
		ok = new(bool)
	}
	return func(err error, req *http.Request, res *http.Response) {
		if res != nil && ErrorKindFrom(err) == KindInvalid {
			*ok = h(res) == nil
		}
	}
}
