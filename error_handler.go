package requests

import (
	"errors"
	"net/http"
)

// ErrorHandler is a function accepted by Builder.OnError.
// Note that the error, request, and response may all be nil
// depending on the error encountered and the effect of prior handlers.
type ErrorHandler = func(ErrorKind, error, *http.Request, *http.Response) error

// ValidatorHandler converts a ResponseHandler into an ErrorHandler for invalid responses.
// The ResponseHandler only runs if the ErrorHandler encounters a validation error.
// If the ResponseHandler succeeds, ErrInvalidHandled is returned.
func ValidatorHandler(h ResponseHandler) ErrorHandler {
	return func(kind ErrorKind, err error, req *http.Request, res *http.Response) error {
		if kind == KindInvalidErr && res != nil {
			if err := h(res); err != nil {
				return err
			}
		}
		return ErrInvalidHandled
	}
}

var ErrInvalidHandled = errors.New("handled recovery from invalid response")
