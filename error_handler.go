package requests

import (
	"errors"
	"net/http"
	"net/url"
)

// ErrorKind indicates where an error was returned in the process of building, validating, and handling a request.
type ErrorKind int8

//go:generate stringer -type=ErrorKind

// Enum values for type ErrorKind
const (
	ErrorKindURL       ErrorKind = iota // error building URL
	ErrorKindRequest                    // error building the request
	ErrorKindConnect                    // error connecting
	ErrorKindValidator                  // validator error
	ErrorKindHandler                    // handler error
)

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

// StatusCode returns Response.StatusCode or 0 if Response is nil
func (ep *OnErrorParams) StatusCode() int {
	if ep.Response == nil {
		return 0
	}
	return ep.Response.StatusCode
}

// ErrorHandler is a function accepted by Builder.OnError.
// Callbacks may modify the fields of ErrorParams.
type ErrorHandler = func(*OnErrorParams)

// ValidatorHandler converts a ResponseHandler into an ErrorHandler for invalid responses.
// The ResponseHandler only runs if the ErrorHandler encounters a validation error.
// If the ResponseHandler succeeds, ErrInvalidHandled is returned.
func ValidatorHandler(h ResponseHandler) ErrorHandler {
	return func(ep *OnErrorParams) {
		if ep.Kind() == ErrorKindValidator && ep.Response != nil {
			if err := h(ep.Response); err == nil { // recovered handler
				ep.Error = ErrInvalidHandled
			}
		}
	}
}

var ErrInvalidHandled = errors.New("handled recovery from invalid response")
