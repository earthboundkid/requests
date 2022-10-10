package requests

import "errors"

// ErrorKind indicates where an error was returned in the process of building, validating, and handling a request.
type ErrorKind int8

//go:generate stringer -type=ErrorKind

// Enum values for type ErrorKind
const (
	KindUnknown    ErrorKind = iota // Not a Builder error
	KindNone                        // error is nil
	KindURLErr                      // error building URL
	KindBodyGetErr                  // error getting request body
	KindMethodErr                   // request method was invalid
	KindContextErr                  // request context was nil
	KindConnectErr                  // error connecting
	KindInvalidErr                  // validation failure
	KindHandlerErr                  // handler error
)

type ek struct {
	kind  ErrorKind
	cause error
}

func (e ek) Error() string {
	return e.cause.Error()
}

func (e ek) Unwrap() error {
	return e.cause
}

// HasKindErr extracts the ErrorKind from an error.
// Nil errors return KindNone.
// Errors not from a Builder return KindUnknown.
func HasKindErr(err error) ErrorKind {
	if err == nil {
		return KindNone
	}
	var e ek
	errors.As(err, &e) // defaults to KindUknown
	return e.kind
}
