package requests

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Builder is a convenient way to build, send, and handle HTTP requests.
// Builder has a fluent API with methods returning a pointer to the same
// struct, which allows for declaratively describing a request by method chaining.
//
// Builder can be thought of as having the following phases:
//
// Set the base URL for a request with requests.URL then customize it with
// Scheme, Host, Hostf, Path, Pathf, and Param.
//
// Set the method for a request with Method or use the Get, Post, and Put
// methods. By default, requests without a body are GET and those with a
// body are POST.
//
// Set headers with Header or set conventional header keys with Accept,
// CacheControl, ContentType, UserAgent, BasicAuth, and Bearer.
//
// Set the http.Client to use for a request with Client.
//
// Set the body of the request, if any, with Body or use built in BodyBytes,
// BodyForm, BodyJSON, BodyReader, or BodyWriter.
//
// Add a response validator to the Builder with AddValidator or use the built
// in CheckStatus, CheckContentType, and Peek.
//
// Set a handler for a response with Handle or use the built in ToJSON,
// ToString, ToBytesBuffer, or ToWriter.
//
// Fetch creates an http.Request with Request and sends it via the underlying
// http.Client with Do.
//
// Config can be used to set several options on a Builder at once.
//
// In many cases, it will be possible to set most options for an API endpoint
// in a Builder at the package or struct level and then call Clone in a
// function to add request specific details for the URL, parameters, headers,
// body, or handler. The zero value of Builder is usable.
type Builder struct {
	baseurl      string
	scheme, host string
	paths        []string
	params       []multimap
	headers      []multimap
	getBody      BodyGetter
	method       string
	cl           *http.Client
	validators   []ResponseHandler
	handler      ResponseHandler
}

type multimap struct {
	key    string
	values []string
}

// URL creates a new Builder suitable for method chaining.
func URL(baseurl string) *Builder {
	var rb Builder
	rb.baseurl = baseurl
	return &rb
}

// Client sets the http.Client to use for requests. If nil, it uses http.DefaultClient.
func (rb *Builder) Client(cl *http.Client) *Builder {
	rb.cl = cl
	return rb
}

// Scheme sets the scheme for a request. It overrides the URL function.
func (rb *Builder) Scheme(scheme string) *Builder {
	rb.scheme = scheme
	return rb
}

// Host sets the host for a request. It overrides the URL function.
func (rb *Builder) Host(host string) *Builder {
	rb.host = host
	return rb
}

// Hostf calls Host with fmt.Sprintf.
func (rb *Builder) Hostf(format string, a ...interface{}) *Builder {
	return rb.Host(fmt.Sprintf(format, a...))
}

// Path joins a path to a request per the path joining rules of RFC 3986.
// If the path begins with /, it overrides any existing path.
// If the path begins with ./ or ../, the final path will be rewritten in its absolute form when creating a request.
func (rb *Builder) Path(path string) *Builder {
	rb.paths = append(rb.paths, path)
	return rb
}

// Pathf calls Path with fmt.Sprintf.
func (rb *Builder) Pathf(format string, a ...interface{}) *Builder {
	return rb.Path(fmt.Sprintf(format, a...))
}

// Param sets a query parameter on a request. It overwrites the existing values of a key.
func (rb *Builder) Param(key string, values ...string) *Builder {
	rb.params = append(rb.params, multimap{key, values})
	return rb
}

// Header sets a header on a request. It overwrites the existing values of a key.
func (rb *Builder) Header(key string, values ...string) *Builder {
	rb.headers = append(rb.headers, multimap{key, values})
	return rb
}

// Accept sets the Accept header for a request.
func (rb *Builder) Accept(contentTypes string) *Builder {
	return rb.Header("Accept", contentTypes)
}

// CacheControl sets the client-side Cache-Control directive for a request.
func (rb *Builder) CacheControl(directive string) *Builder {
	return rb.Header("Cache-Control", directive)
}

// ContentType sets the Content-Type header on a request.
func (rb *Builder) ContentType(ct string) *Builder {
	return rb.Header("Content-Type", ct)
}

// UserAgent sets the User-Agent header.
func (rb *Builder) UserAgent(s string) *Builder {
	return rb.Header("User-Agent", s)
}

// BasicAuth sets the Authorization header to a basic auth credential.
func (rb *Builder) BasicAuth(username, password string) *Builder {
	auth := username + ":" + password
	v := base64.StdEncoding.EncodeToString([]byte(auth))
	return rb.Header("Authorization", "Basic "+v)
}

// Bearer sets the Authorization header to a bearer token.
func (rb *Builder) Bearer(token string) *Builder {
	return rb.Header("Authorization", "Bearer "+token)
}

// Method sets the HTTP method for a request.
func (rb *Builder) Method(method string) *Builder {
	rb.method = method
	return rb
}

// Get sets HTTP method to GET.
func (rb *Builder) Get() *Builder {
	return rb.Method(http.MethodGet)
}

// Head sets HTTP method to HEAD.
func (rb *Builder) Head() *Builder {
	return rb.Method(http.MethodHead)
}

// Post sets HTTP method to POST.
func (rb *Builder) Post() *Builder {
	return rb.Method(http.MethodPost)
}

// Put sets HTTP method to PUT.
func (rb *Builder) Put() *Builder {
	return rb.Method(http.MethodPut)
}

// Body sets the BodyGetter to use to build the body of a request.
// The provided BodyGetter is used as an http.Request.GetBody func.
// It implicitly sets method to POST.
func (rb *Builder) Body(src BodyGetter) *Builder {
	rb.getBody = src
	return rb
}

// BodyReader sets the Builder's request body to r.
func (rb *Builder) BodyReader(r io.Reader) *Builder {
	return rb.Body(BodyReader(r))
}

// BodyWriter pipes writes from w to the Builder's request body.
func (rb *Builder) BodyWriter(f func(w io.Writer) error) *Builder {
	return rb.Body(BodyWriter(f))
}

// BodyBytes sets the Builder's request body to b.
func (rb *Builder) BodyBytes(b []byte) *Builder {
	return rb.Body(BodyBytes(b))
}

// BodyJSON sets the Builder's request body to the marshaled JSON.
// It also sets ContentType to "application/json".
func (rb *Builder) BodyJSON(v interface{}) *Builder {
	return rb.
		Body(BodyJSON(v)).
		ContentType("application/json")
}

// BodyForm sets the Builder's request body to the encoded form.
// It also sets the ContentType to "application/x-www-form-urlencoded".
func (rb *Builder) BodyForm(data url.Values) *Builder {
	return rb.
		Body(BodyForm(data)).
		ContentType("application/x-www-form-urlencoded")
}

// AddValidator adds a response validator to the Builder.
// Adding a validator disables DefaultValidator.
// To disable all validation, just add nil.
func (rb *Builder) AddValidator(h ResponseHandler) *Builder {
	rb.validators = append(rb.validators, h)
	return rb
}

// CheckStatus adds a validator for status code of a response.
func (rb *Builder) CheckStatus(acceptStatuses ...int) *Builder {
	return rb.AddValidator(CheckStatus(acceptStatuses...))
}

// CheckContentType adds a validator for the content type header of a response.
func (rb *Builder) CheckContentType(cts ...string) *Builder {
	return rb.AddValidator(CheckContentType(cts...))
}

// CheckPeek adds a validator that peeks at the first n bytes of a response body.
func (rb *Builder) CheckPeek(n int, f func([]byte) error) *Builder {
	return rb.AddValidator(CheckPeek(n, f))
}

// Handle sets the response handler for a Builder.
// To use multiple handlers, use ChainHandlers.
func (rb *Builder) Handle(h ResponseHandler) *Builder {
	rb.handler = h
	return rb
}

// ToJSON sets the Builder to decode a response as a JSON object
func (rb *Builder) ToJSON(v interface{}) *Builder {
	return rb.Handle(ToJSON(v))
}

// ToString sets the Builder to write the response body to the provided string pointer.
func (rb *Builder) ToString(sp *string) *Builder {
	return rb.Handle(ToString(sp))
}

// ToBytesBuffer sets the Builder to write the response body to the provided bytes.Buffer.
func (rb *Builder) ToBytesBuffer(buf *bytes.Buffer) *Builder {
	return rb.Handle(ToBytesBuffer(buf))
}

// ToWriter sets the Builder to copy the response body into w.
func (rb *Builder) ToWriter(w io.Writer) *Builder {
	return rb.Handle(ToWriter(w))
}

// Config allows Builder to be extended by functions that set several options at once.
func (rb *Builder) Config(cfg Config) *Builder {
	cfg(rb)
	return rb
}

// Clone creates a new Builder suitable for independent mutation.
func (rb *Builder) Clone() *Builder {
	rb2 := *rb
	rb2.paths = rb2.paths[0:len(rb2.paths):len(rb2.paths)]
	rb2.headers = rb2.headers[0:len(rb2.headers):len(rb2.headers)]
	rb2.params = rb2.params[0:len(rb2.params):len(rb2.params)]
	rb2.validators = rb2.validators[0:len(rb2.validators):len(rb2.validators)]
	return &rb2
}

// Request builds a new http.Request with its context set.
func (rb *Builder) Request(ctx context.Context) (req *http.Request, err error) {
	u, err := url.Parse(rb.baseurl)
	if err != nil {
		return nil, fmt.Errorf("could not initialize with base URL %q: %w", u, err)
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	if rb.scheme != "" {
		u.Scheme = rb.scheme
	}
	if rb.host != "" {
		u.Host = rb.host
	}
	for _, p := range rb.paths {
		u.Path = u.ResolveReference(&url.URL{Path: p}).Path
	}
	if len(rb.params) > 0 {
		q := u.Query()
		for _, kv := range rb.params {
			q[kv.key] = kv.values
		}
		u.RawQuery = q.Encode()
	}
	var body io.ReadCloser
	if rb.getBody != nil {
		if body, err = rb.getBody(); err != nil {
			return nil, err
		}
	}
	method := http.MethodGet
	if rb.getBody != nil {
		method = http.MethodPost
	}
	if rb.method != "" {
		method = rb.method
	}
	req, err = http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.GetBody = rb.getBody

	for _, kv := range rb.headers {
		req.Header[http.CanonicalHeaderKey(kv.key)] = kv.values
	}
	return req, nil
}

// Do calls the underlying http.Client and validates and handles any resulting response. The response body is closed after all validators and the handler run.
func (rb *Builder) Do(req *http.Request) (err error) {
	cl := http.DefaultClient
	if rb.cl != nil {
		cl = rb.cl
	}
	res, err := cl.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	validators := rb.validators
	if len(validators) == 0 {
		validators = []ResponseHandler{DefaultValidator}
	}
	if err = ChainHandlers(validators...)(res); err != nil {
		return err
	}
	h := consumeBody
	if rb.handler != nil {
		h = rb.handler
	}
	if err = h(res); err != nil {
		return err
	}
	return nil
}

// Fetch builds a request, sends it, and handles the response.
func (rb *Builder) Fetch(ctx context.Context) (err error) {
	req, err := rb.Request(ctx)
	if err != nil {
		return err
	}
	return rb.Do(req)
}
