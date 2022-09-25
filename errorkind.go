package requests

import "errors"

// ErrorKind indicates where an error was returned in the process of building, validating, and handling a request.
type ErrorKind int

//go:generate stringer -type=ErrorKind

const (
	ErrorKindNone ErrorKind = iota
	ErrorKindUnknown
	ErrorKindURLParse
	ErrorKindBodyGet
	ErrorKindUnknownMethod
	ErrorKindNilContext
	ErrorKindConnection
	ErrorKindValidator
	ErrorKindHandler
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

// ErrorKindFrom extracts the ErrorKind from an error.
// Nil errors return ErrorKindNone. Errors not passed through
// a Builder return ErrorKindUnknown.
func ErrorKindFrom(err error) ErrorKind {
	if err == nil {
		return ErrorKindNone
	}
	var e ek
	if !errors.As(err, &e) {
		return ErrorKindUnknown
	}
	return e.kind
}
