package requests

// ErrorKind indicates where an error was returned in the process of building, validating, and handling a request.
// Errors returned by Builder can be tested for their ErrorKind using errors.Is or errors.As.
type ErrorKind int8

//go:generate stringer -type=ErrorKind

// Enum values for type ErrorKind
const (
	ErrURL       ErrorKind = iota // error building URL
	ErrRequest                    // error building the request
	ErrTransport                  // error connecting
	ErrValidator                  // validator error
	ErrHandler                    // handler error
)

func (ek ErrorKind) Error() string {
	return ek.String()
}

type ekwrapper struct {
	kind ErrorKind
	error
}

func (ekw ekwrapper) Is(target error) bool {
	return ekw.kind == target
}

func (ekw ekwrapper) As(target any) bool {
	if ekp, ok := target.(*ErrorKind); ok {
		*ekp = ekw.kind
		return true
	}
	return false
}

func (ekw ekwrapper) Unwrap() error {
	return ekw.error
}
