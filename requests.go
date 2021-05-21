package requests

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

type Builder struct {
	cl         *http.Client
	host, path string
	params     [][2]string
	headers    [][2]string
	url        *url.URL
	err        error
	method     string
	body       BodySource
	validator  ResponseHandler
	handler    ResponseHandler
}

// URL creates a new Builder suitable for method chaining.
func URL(u string) *Builder {
	var rb Builder
	rb.url, rb.err = url.Parse(u)
	if rb.err != nil {
		rb.err = fmt.Errorf("could not initialize with URL %q: %w", u, rb.err)
	}
	return &rb
}

// Client sets the http.Client to use for requests. If nil, it uses http.DefaultClient.
func (rb *Builder) Client(cl *http.Client) *Builder {
	rb.cl = cl
	return rb
}

// Host sets the host for a request. It overrides the URL function.
func (rb *Builder) Host(host string) *Builder {
	rb.host = host
	return rb
}

// Path sets the path for a request. It overrides the URL function.
func (rb *Builder) Path(path string) *Builder {
	rb.path = path
	return rb
}

// Param adds a query parameter to a request. It does not remove existing parameters.
func (rb *Builder) Param(key, value string) *Builder {
	rb.params = append(rb.params, [2]string{key, value})
	return rb
}

// Method sets the HTTP method for a request.
func (rb *Builder) Method(method string) *Builder {
	rb.method = method
	return rb
}

func (rb *Builder) Get() *Builder {
	return rb.Method(http.MethodGet)
}

func (rb *Builder) Post() *Builder {
	return rb.Method(http.MethodPost)
}

func (rb *Builder) Put() *Builder {
	return rb.Method(http.MethodPut)
}

// BodySource provides a builder with a source for a request body.
type BodySource = func() (io.ReadCloser, string, error)

// Body sets the BodySource for a request. It implicitly sets method to POST.
func (rb *Builder) Body(src BodySource) *Builder {
	rb.body = src
	return rb
}

func Reader(r io.Reader, contentType string) BodySource {
	return func() (io.ReadCloser, string, error) {
		if rc, ok := r.(io.ReadCloser); ok {
			return rc, contentType, nil
		}
		return io.NopCloser(r), contentType, nil
	}
}

func (rb *Builder) Reader(r io.Reader, contentType string) *Builder {
	return rb.Body(Reader(r, contentType))
}

func Bytes(b []byte, contentType string) BodySource {
	return func() (io.ReadCloser, string, error) {
		return io.NopCloser(bytes.NewReader(b)), contentType, nil
	}
}

func (rb *Builder) Bytes(b []byte, contentType string) *Builder {
	return rb.Body(Bytes(b, contentType))
}

func JSONMarshal(v interface{}) BodySource {
	return func() (r io.ReadCloser, contentType string, err error) {
		contentType = "application/json"
		b, err := json.Marshal(v)
		if err != nil {
			return
		}
		r = io.NopCloser(bytes.NewReader(b))
		return
	}
}

func (rb *Builder) JSONMarshal(v interface{}) *Builder {
	return rb.Body(JSONMarshal(v))
}

func Form(data url.Values) BodySource {
	return func() (r io.ReadCloser, contentType string, err error) {
		return io.NopCloser(strings.NewReader(data.Encode())),
			"application/x-www-form-urlencoded", nil
	}
}

func (rb *Builder) Form(data url.Values) *Builder {
	return rb.Body(Form(data))
}

// ResponseHandler is used to validate or handle the response to a request.
type ResponseHandler = func(*http.Response) error

// ChainHandlers allows for the composing of validators or response handlers.
func ChainHandlers(handlers ...ResponseHandler) ResponseHandler {
	return func(r *http.Response) error {
		for _, h := range handlers {
			if h == nil {
				continue
			}
			if err := h(r); err != nil {
				return err
			}
		}
		return nil
	}
}

func (rb *Builder) Validate(h ResponseHandler) *Builder {
	rb.validator = h
	return rb
}

// CheckStatus validates the response has an acceptable status code.
func CheckStatus(acceptStatuses ...int) ResponseHandler {
	return func(resp *http.Response) error {
		for _, code := range acceptStatuses {
			if resp.StatusCode == code {
				return nil
			}
		}

		return StatusError{
			resp.Request.URL.Redacted(),
			resp.Status,
			resp.StatusCode,
		}
	}
}

func (rb *Builder) CheckStatus(acceptStatuses ...int) *Builder {
	return rb.Validate(CheckStatus(acceptStatuses...))
}

// DefaultValidator is the status check applied by Builder unless otherwise specified.
var DefaultValidator ResponseHandler = CheckStatus(
	http.StatusOK,
	http.StatusCreated,
	http.StatusAccepted,
	http.StatusNonAuthoritativeInfo,
	http.StatusNoContent,
)

// StatusError is the error type produced by CheckStatus.
type StatusError struct {
	URL, Status string
	StatusCode  int
}

// Error fulfills the error interface.
func (se StatusError) Error() string {
	return fmt.Sprintf("unexpected status for %s: %s",
		se.URL, se.Status)
}

// HasStatusErr returns true if err is a StatusError caused by any of the code given.
func HasStatusErr(err error, codes ...int) bool {
	if se := (StatusError{}); errors.As(err, &se) {
		for _, code := range codes {
			if se.StatusCode == code {
				return true
			}
		}
	}
	return false
}

// MatchContentType validates that a response has the given content type.
func MatchContentType(ct string) ResponseHandler {
	return func(resp *http.Response) error {
		mt, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
		if err != nil {
			return fmt.Errorf("problem matching Content-Type: %w", err)
		}
		if mt != ct {
			return fmt.Errorf("unexpected Content-Type: %s", mt)
		}
		return nil
	}
}

func (rb *Builder) Handle(h ResponseHandler) *Builder {
	rb.handler = h
	return rb
}

func consumeBody(res *http.Response) (err error) {
	const maxDiscardSize = 640 * 1 << 10
	if _, err = io.CopyN(io.Discard, res.Body, maxDiscardSize); err == io.EOF {
		err = nil
	}
	return err
}

func JSONUnmarshal(v interface{}) ResponseHandler {
	return func(res *http.Response) error {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(data, v); err != nil {
			return err
		}
		return nil
	}
}

func (rb *Builder) JSONUnmarshal(v interface{}) *Builder {
	return rb.Handle(JSONUnmarshal(v))
}

func Buffer(buf *bytes.Buffer) ResponseHandler {
	return func(res *http.Response) error {
		_, err := io.Copy(buf, res.Body)
		return err
	}
}

func (rb *Builder) Buffer(buf *bytes.Buffer) *Builder {
	return rb.Handle(Buffer(buf))
}

func BufioReader(f func(r *bufio.Reader) error) ResponseHandler {
	return func(res *http.Response) error {
		return f(bufio.NewReader(res.Body))
	}
}

func (rb *Builder) BufioReader(f func(r *bufio.Reader) error) *Builder {
	return rb.Handle(BufioReader(f))
}

// Clone creates a new Builder suitable for independent mutation.
func (rb *Builder) Clone() *Builder {
	rb2 := *rb
	rb2.headers = rb2.headers[0:len(rb2.headers):len(rb2.headers)]
	rb2.params = rb2.params[0:len(rb2.params):len(rb2.params)]
	u := *rb.url
	rb2.url = &u
	return &rb2
}

// Request builds a new http.Request with its context set.
func (rb *Builder) Request(ctx context.Context) (req *http.Request, err error) {
	if rb.err != nil {
		return nil, err
	}
	method := http.MethodGet
	if rb.body != nil {
		method = http.MethodPost
	}
	if rb.method != "" {
		method = rb.method
	}
	if rb.url == nil {
		if rb.host == "" {
			return nil, fmt.Errorf("must set a URL to connect to")
		}
		rb.url = &url.URL{}
	}
	if rb.url.Scheme == "" {
		rb.url.Scheme = "https"
	}
	if rb.host != "" {
		rb.url.Host = rb.host
	}
	if rb.path != "" {
		rb.url.Path = rb.path
	}
	if len(rb.params) > 0 {
		q := rb.url.Query()
		for _, kv := range rb.params {
			q.Add(kv[0], kv[1])
		}
		rb.url.RawQuery = q.Encode()
	}
	u := rb.url.String()
	var body io.ReadCloser
	var ct string
	if rb.body != nil {
		if body, ct, err = rb.body(); err != nil {
			return nil, err
		}
	}
	req, err = http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, err
	}
	for _, pair := range rb.headers {
		req.Header.Add(pair[0], pair[1])
	}
	if req.Header.Get("Content-Type") == "" && ct != "" {
		req.Header.Set("Content-Type", ct)
	}
	return req, nil
}

// Do calls the underlying http.Client and handles any resulting response.
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

	h := DefaultValidator
	if rb.validator != nil {
		h = rb.validator
	}
	if err = h(res); err != nil {
		return err
	}
	h = consumeBody
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
