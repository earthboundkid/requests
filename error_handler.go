package requests

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// ErrorKind indicates where an error was returned in the process of building, validating, and handling a request.
// Errors returned by Builder can be tested for their ErrorKind using errors.Is or errors.As.
type ErrorKind int8

//go:generate stringer -type=ErrorKind

// Enum values for type ErrorKind
const (
	ErrUnknown   ErrorKind = iota // error not from requests.Builder
	ErrURL                        // error building URL
	ErrRequest                    // error building the request
	ErrConnect                    // error connecting
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

// OnErrorParams is a struct used by ErrorHandlers to describe an error encounted by a Builder.
// Note that Error, Request, and Response may all be nil
// depending on the error encountered and the effect of prior handlers.
type OnErrorParams struct {
	Error    error
	Request  *http.Request
	Response *http.Response
	kind     ErrorKind
	rb       *Builder
}

// Kind returns ErrorKind that created the OnErrorParams.
// It will stay the same even if OnErrorParams.Error is changed.
func (ep *OnErrorParams) Kind() ErrorKind {
	return ep.kind
}

// Status returns Response.Status or "" if Response is nil
func (ep *OnErrorParams) Status() string {
	if ep.Response == nil {
		return ""
	}
	return ep.Response.Status
}

// Method the HTTP Method of the Builder originating the OnErrorParams.
// Note also that Response.Request.Method and Request.Method may differ
// if a request has been redirected or altered by a Transport.
func (ep *OnErrorParams) Method() string {
	return ep.rb.getMethod()
}

// URL calls URL() on the Builder originating the OnErrorParams.
// Note also that Response.Request.URL and Request.URL may differ
// if a request has been redirected or altered by a Transport.
func (ep *OnErrorParams) URL() *url.URL {
	u, _ := ep.rb.URL()
	return u
}

// StatusCode returns Response.StatusCode or 0 if Response is nil.
func (ep *OnErrorParams) StatusCode() int {
	if ep.Response == nil {
		return 0
	}
	return ep.Response.StatusCode
}

// Context returns Request.Context() or context.Background if Request is nil.
func (ep *OnErrorParams) Context() context.Context {
	if ep.Request == nil {
		return context.Background()
	}
	return ep.Request.Context()
}

// ErrorHandler is a function accepted by Builder.OnError.
// Callbacks may modify the fields of ErrorParams.
type ErrorHandler = func(*OnErrorParams)

// ValidatorHandler converts a ResponseHandler into an ErrorHandler for invalid responses.
// The ResponseHandler only runs if the ErrorHandler encounters a validation error.
// If the ResponseHandler succeeds, ErrInvalidHandled is returned.
func ValidatorHandler(h ResponseHandler) ErrorHandler {
	return func(ep *OnErrorParams) {
		if ep.Kind() == ErrValidator && ep.Response != nil {
			if err := h(ep.Response); err == nil { // recovered handler
				ep.Error = ErrInvalidHandled
			}
		}
	}
}

var ErrInvalidHandled = errors.New("handled recovery from invalid response")
