package requests

// ErrorKind indicates where an error was returned in the process of building, validating, and handling a request.
type ErrorKind int8

//go:generate stringer -type=ErrorKind

// Enum values for type ErrorKind
const (
	ErrorKindUnknown   ErrorKind = iota // Not a Builder error
	ErrorKindURL                        // error building URL
	ErrorKindBodyGet                    // error getting request body
	ErrorKindMethod                     // request method was invalid
	ErrorKindContext                    // request context was nil
	ErrorKindConnect                    // error connecting
	ErrorKindValidator                  // validator error
	ErrorKindHandler                    // handler error
)

// ErrorKindError is an error that can return its underlying ErrorKind.
// Errors returned by Builder conform to ErrorKindError.
type ErrorKindError interface {
	error
	Kind() ErrorKind
}

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

func (e ek) Kind() ErrorKind {
	return e.kind
}
