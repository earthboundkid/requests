package requests

import (
	"net/http"
)

type ErrorKind int

type ErrorHandler func(kind ErrorKind, cause error, resp *http.Response) error

const (
	ErrorKindAny ErrorKind = iota
	ErrorKindConnectionErr
	ErrorKindValidationErr
	ErrorKindHandlerErr
)

func (rb *Builder) handleError(kind ErrorKind, cause error, resp *http.Response) error {
	// If we've got any "any" handler, always use that one
	if anyhandler, ok := rb.errorHandlers[ErrorKindAny]; ok {
		return anyhandler(kind, cause, resp)
	}

	// Otherwise lookup the specific kind
	specificHandler, ok := rb.errorHandlers[kind]
	// If we don't have one, then just return the original error unchanged
	if !ok {
		return cause
	}

	return specificHandler(kind, cause, resp)
}
