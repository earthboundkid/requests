package reqtest

import (
	"bufio"
	"bytes"
	"io/fs"
	"net/http"
	"os"

	"github.com/carlmjohnson/requests"
)

// ReplayString returns an http.RoundTripper that always responds with a
// request built from rawResponse. It is intended for use in one-off tests.
func ReplayString(rawResponse string) requests.Transport {
	return requests.ReplayString(rawResponse)
}

// Record returns an http.RoundTripper that writes out its
// requests and their responses to text files in basepath.
// Requests are named according to a hash of their contents.
// Responses are named according to the request that made them.
func Record(rt http.RoundTripper, basepath string) requests.Transport {
	return requests.Record(rt, basepath)
}

// Replay returns an http.RoundTripper that reads its
// responses from text files in basepath.
// Responses are looked up according to a hash of the request.
// Response file names may optionally be prefixed with comments for better human organization.
func Replay(basepath string) requests.Transport {
	return requests.Replay(basepath)
}

// ReplayFile returns an http.RoundTripper that reads its
// response from the response file at path.
func ReplayFile(path string) requests.Transport {
	return requests.RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		r := bufio.NewReader(bytes.NewReader(b))
		return http.ReadResponse(r, req)
	})
}

// ReplayFS returns an http.RoundTripper that reads its
// responses from text files in the fs.FS.
// Responses are looked up according to a hash of the request.
// Response file names may optionally be prefixed with comments for better human organization.
func ReplayFS(fsys fs.FS) requests.Transport {
	return requests.ReplayFS(fsys)
}

// Caching returns an http.RoundTripper that attempts to read its
// responses from text files in basepath. If the response is absent,
// it caches the result of issuing the request with rt in basepath.
// Requests are named according to a hash of their contents.
// Responses are named according to the request that made them.
func Caching(rt http.RoundTripper, basepath string) requests.Transport {
	return requests.Caching(rt, basepath)
}
