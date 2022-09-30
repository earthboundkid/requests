package requests

import "errors"

// ErrorKind indicates where an error was returned in the process of building, validating, and handling a request.
type ErrorKind int8

//go:generate stringer -type=ErrorKind

// Enum values for type ErrorKind
const (
	KindNone ErrorKind = iota
	KindUnknown
	KindURLErr
	KindBodyGet
	KindBadMethod
	KindNilContext
	KindConnectErr
	KindInvalid
	KindHandlerErr
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
// Nil errors return KindNone.
// Errors not from a Builder return KindUnknown.
func ErrorKindFrom(err error) ErrorKind {
	if err == nil {
		return KindNone
	}
	var e ek
	if !errors.As(err, &e) {
		return KindUnknown
	}
	return e.kind
}
