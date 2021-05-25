package requests

import (
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

// AddCookieJar adds the standard public suffix list cookiejar to an http.Client.
// Because it modifies the client in place, cl should not be nil.
func AddCookieJar(cl *http.Client) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	// As of Go 1.16, cookiejar.New err is hardcoded nil
	if err != nil {
		panic(err)
	}
	cl.Jar = jar
}
