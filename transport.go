package requests

import (
	"bufio"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Transport is an alias of http.RoundTripper for documentation purposes.
type Transport = http.RoundTripper

// RoundTripFunc is an adaptor to use a function as an http.RoundTripper.
type RoundTripFunc func(req *http.Request) (res *http.Response, err error)

// RoundTrip implements http.RoundTripper.
func (rtf RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return rtf(r)
}

var _ Transport = RoundTripFunc(nil)

// ReplayString returns an http.RoundTripper that always responds with a
// request built from rawResponse. It is intended for use in one-off tests.
func ReplayString(rawResponse string) Transport {
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		r := bufio.NewReader(strings.NewReader(rawResponse))
		res, err = http.ReadResponse(r, req)
		return
	})
}

// UserAgentTransport returns a wrapped http.RoundTripper that sets the User-Agent header on requests to s.
func UserAgentTransport(rt http.RoundTripper, s string) Transport {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		r2 := *req
		r2.Header = r2.Header.Clone()
		r2.Header.Set("User-Agent", s)
		return rt.RoundTrip(&r2)
	})
}

// PermitURLTransport returns a wrapped http.RoundTripper that rejects any requests whose URL doesn't match the provided regular expression string.
//
// PermitURLTransport will panic if the regexp does not compile.
func PermitURLTransport(rt http.RoundTripper, regex string) Transport {
	if rt == nil {
		rt = http.DefaultTransport
	}
	re := regexp.MustCompile(regex)
	reErr := fmt.Errorf("requested URL not permitted by regexp: %s", regex)
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		if u := req.URL.String(); !re.MatchString(u) {
			return nil, reErr
		}
		return rt.RoundTrip(req)
	})
}
