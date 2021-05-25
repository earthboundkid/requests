package requests

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/carlmjohnson/crockford"
)

// Replay returns an http.RoundTripper that reads its
// responses from text files in basepath.
// Responses are looked up according to a hash of the request.
func Replay(basepath string) http.RoundTripper {
	return ReplayFS(os.DirFS(basepath))
}

// ReplayFS returns an http.RoundTripper that reads its
// responses from text files in the fs.FS.
// Responses are looked up according to a hash of the request.
func ReplayFS(fsys fs.FS) http.RoundTripper {
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("problem while replaying transport: %w", err)
			}
		}()
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}
		md5 := crockford.MD5(crockford.Lower, b)
		b, err = fs.ReadFile(fsys, md5[:5]+".res.txt")
		if err != nil {
			return nil, err
		}
		r := bufio.NewReader(bytes.NewReader(b))
		res, err = http.ReadResponse(r, req)
		return
	})
}
