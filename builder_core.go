package requests

import (
	"context"
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
// Set the method for a request with Method or use the Delete, Patch, and Put
// methods. By default, requests without a body are GET and those with a
// body are POST.
//
// Set headers with Header or set conventional header keys with Accept,
// CacheControl, ContentType, Cookie, UserAgent, BasicAuth, and Bearer.
//
// Set the http.Client to use for a request with Client and/or set an
// http.RoundTripper with Transport.
//
// Set the body of the request, if any, with Body or use built in BodyBytes,
// BodyFile, BodyForm, BodyJSON, BodyReader, or BodyWriter.
//
// Add a response validator to the Builder with AddValidator or use the built
// in CheckStatus, CheckContentType, CopyHeaders, and Peek.
//
// Set a handler for a response with Handle or use the built in ToHeaders,
// ToJSON, ToString, ToBytesBuffer, or ToWriter.
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
	cookies      []kvpair
	getBody      BodyGetter
	method       string
	cl           *http.Client
	rt           http.RoundTripper
	validators   []ResponseHandler
	handler      ResponseHandler
	onerrs       []ErrorHandler
}

type multimap struct {
	key    string
	values []string
}

type kvpair struct {
	key, value string
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

// Transport sets the http.RoundTripper to use for requests.
// If set, it makes a shallow copy of the http.Client before modifying it.
func (rb *Builder) Transport(rt http.RoundTripper) *Builder {
	rb.rt = rt
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

// Path joins a path to a request per the path joining rules of RFC 3986.
// If the path begins with /, it overrides any existing path.
// If the path begins with ./ or ../, the final path will be rewritten in its absolute form when creating a request.
func (rb *Builder) Path(path string) *Builder {
	rb.paths = append(rb.paths, path)
	return rb
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

// Cookie adds a cookie to a request.
// Unlike other headers, adding a cookie does not overwrite existing values.
func (rb *Builder) Cookie(name, value string) *Builder {
	rb.cookies = append(rb.cookies, kvpair{name, value})
	return rb
}

// Method sets the HTTP method for a request.
func (rb *Builder) Method(method string) *Builder {
	rb.method = method
	return rb
}

// Body sets the BodyGetter to use to build the body of a request.
// The provided BodyGetter is used as an http.Request.GetBody func.
// It implicitly sets method to POST.
func (rb *Builder) Body(src BodyGetter) *Builder {
	rb.getBody = src
	return rb
}

// AddValidator adds a response validator to the Builder.
// Adding a validator disables DefaultValidator.
// To disable all validation, just add nil.
func (rb *Builder) AddValidator(h ResponseHandler) *Builder {
	rb.validators = append(rb.validators, h)
	return rb
}

// Handle sets the response handler for a Builder.
// To use multiple handlers, use ChainHandlers.
func (rb *Builder) Handle(h ResponseHandler) *Builder {
	rb.handler = h
	return rb
}

// Config allows Builder to be extended by functions that set several options at once.
func (rb *Builder) Config(cfgs ...Config) *Builder {
	for _, cfg := range cfgs {
		cfg(rb)
	}
	return rb
}

// OnError adds an ErrorHandler to run if any part of building, validating, or handling a request fails.
func (rb *Builder) OnError(h ErrorHandler) *Builder {
	rb.onerrs = append(rb.onerrs, h)
	return rb
}

func clip[T any](sp *[]T) {
	s := *sp
	*sp = s[:len(s):len(s)]
}

// Clone creates a new Builder suitable for independent mutation.
func (rb *Builder) Clone() *Builder {
	rb2 := *rb
	clip(&rb2.paths)
	clip(&rb2.headers)
	clip(&rb2.params)
	clip(&rb2.cookies)
	clip(&rb2.validators)
	clip(&rb2.onerrs)
	return &rb2
}

// Request builds a new http.Request with its context set.
func (rb *Builder) Request(ctx context.Context) (req *http.Request, err error) {
	u, err := url.Parse(rb.baseurl)
	if err != nil {
		err = fmt.Errorf("could not initialize with base URL %q: %w", u, err)
		err = ek{ErrorKindURLParse, err}
		rb.handleErr(err, nil, nil)
		return nil, err
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
	var body io.Reader
	if rb.getBody != nil {
		if body, err = rb.getBody(); err != nil {
			err = ek{ErrorKindBodyGet, err}
			rb.handleErr(err, nil, nil)
			return nil, err
		}
		if nopper, ok := body.(nopCloser); ok {
			body = nopper.Reader
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
		if ctx == nil {
			err = ek{ErrorKindNilContext, err}
		} else {
			err = ek{ErrorKindUnknownMethod, err}
		}
		rb.handleErr(err, nil, nil)
		return nil, err
	}
	req.GetBody = rb.getBody

	for _, kv := range rb.headers {
		req.Header[http.CanonicalHeaderKey(kv.key)] = kv.values
	}
	for _, kv := range rb.cookies {
		req.AddCookie(&http.Cookie{
			Name:  kv.key,
			Value: kv.value,
		})
	}
	return req, nil
}

// Do calls the underlying http.Client and validates and handles any resulting response. The response body is closed after all validators and the handler run.
func (rb *Builder) Do(req *http.Request) (err error) {
	cl := http.DefaultClient
	if rb.cl != nil {
		cl = rb.cl
	}
	if rb.rt != nil {
		cl2 := *cl
		cl2.Transport = rb.rt
		cl = &cl2
	}
	res, err := cl.Do(req)
	if err != nil {
		err = ek{ErrorKindConnection, err}
		rb.handleErr(err, req, nil)
		return err
	}
	defer res.Body.Close()

	validators := rb.validators
	if len(validators) == 0 {
		validators = []ResponseHandler{DefaultValidator}
	}
	if err = ChainHandlers(validators...)(res); err != nil {
		err = ek{ErrorKindValidator, err}
		rb.handleErr(err, req, res)
		return err
	}
	h := consumeBody
	if rb.handler != nil {
		h = rb.handler
	}
	if err = h(res); err != nil {
		err = ek{ErrorKindHandler, err}
		rb.handleErr(err, req, res)
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

func (rb *Builder) handleErr(err error, req *http.Request, res *http.Response) {
	for _, h := range rb.onerrs {
		h(err, req, res)
	}
}
