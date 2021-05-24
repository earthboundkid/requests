package requests

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"

	"github.com/carlmjohnson/crockford"
)

// Record returns an http.RoundTripper that reads its
// responses from text files in basepath.
func Replay(basepath string) http.RoundTripper {
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("problem while replaying transport: %w", err)
			}
		}()
		_ = os.MkdirAll(basepath, 0755)
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}
		md5 := crockford.MD5(crockford.Lower, b)
		name := filepath.Join(basepath, md5[:5]+".res.txt")
		b, err = os.ReadFile(name)
		if err != nil {
			return nil, err
		}
		r := bufio.NewReader(bytes.NewReader(b))
		res, err = http.ReadResponse(r, req)
		return
	})
}
