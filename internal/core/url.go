// Package core handles the core functionality of requests
// apart from any convenience functions and patterns.
package core

import (
	"net/url"

	"github.com/carlmjohnson/requests/internal/minitrue"
	"github.com/carlmjohnson/requests/internal/slicex"
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
	slicex.Clip(&ub2.paths)
	slicex.Clip(&ub2.params)
	return &ub2
}

func (ub *URLBuilder) URL() (u *url.URL, err error) {
	u, err = url.Parse(ub.baseurl)
	if err != nil {
		return new(url.URL), err
	}
	u.Scheme = minitrue.First(ub.scheme, minitrue.First(u.Scheme, "https"))
	u.Host = minitrue.First(ub.host, u.Host)
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
