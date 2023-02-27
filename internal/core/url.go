// Package core handles the core functionality of requests
// apart from any convenience functions and patterns.
package core

import (
	"io"
	"net/url"

	"github.com/carlmjohnson/requests/internal/util"
)

type multimap struct {
	key    string
	values []string
}

type kvpair struct {
	key, value string
}

type URLBuilder struct {
	baseurl      string
	scheme, host string
	paths        []string
	params       []multimap
}

func (ub *URLBuilder) BaseURL(baseurl string) {
	ub.baseurl = baseurl
}

func (ub *URLBuilder) Scheme(scheme string) {
	ub.scheme = scheme
}

func (ub *URLBuilder) Host(host string) {
	ub.host = host

}

func (ub *URLBuilder) Path(path string) {
	ub.paths = append(ub.paths, path)
}

func (ub *URLBuilder) Param(key string, values ...string) {
	ub.params = append(ub.params, multimap{key, values})
}

func (ub *URLBuilder) Clone() *URLBuilder {
	ub2 := *ub
	util.Clip(&ub2.paths)
	util.Clip(&ub2.params)
	return &ub2
}

func (ub *URLBuilder) URL() (u *url.URL, err error) {
	u, err = url.Parse(ub.baseurl)
	if err != nil {
		return new(url.URL), err
	}
	u.Scheme = util.First(ub.scheme, util.First(u.Scheme, "https"))
	u.Host = util.First(ub.host, u.Host)
	for _, p := range ub.paths {
		u.Path = u.ResolveReference(&url.URL{Path: p}).Path
	}
	if len(ub.params) > 0 {
		q := u.Query()
		for _, kv := range ub.params {
			q[kv.key] = kv.values
		}
		u.RawQuery = q.Encode()
	}
	// Reparsing, in case the path rewriting broke the URL
	u, err = url.Parse(u.String())
	if err != nil {
		return new(url.URL), err
	}
	return u, nil
}

type BodyGetter = func() (io.ReadCloser, error)

type RequestBuilder struct {
	headers []multimap
	cookies []kvpair
	getBody BodyGetter
	method  string
}

func (rb *RequestBuilder) Header(key string, values ...string) {
	rb.headers = append(rb.headers, multimap{key, values})
}

func (rb *RequestBuilder) Cookie(name, value string) {
	rb.cookies = append(rb.cookies, kvpair{name, value})
}

func (rb *RequestBuilder) Method(method string) {
	rb.method = method
}

func (rb *RequestBuilder) Body(src BodyGetter) {
	rb.getBody = src
}

// Clone creates a new Builder suitable for independent mutation.
func (rb *RequestBuilder) Clone() *RequestBuilder {
	rb2 := *rb
	util.Clip(&rb2.headers)
	util.Clip(&rb2.cookies)
	return &rb2
}
