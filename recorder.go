package requests

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"

	"github.com/carlmjohnson/crockford"
)

// Record returns an http.RoundTripper that writes out its
// requests and responses made to text files in basepath.
func Record(rt http.RoundTripper, basepath string) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("problem while recording transport: %w", err)
			}
		}()
		_ = os.MkdirAll(basepath, 0755)
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}
		md5 := crockford.MD5(crockford.Lower, b)
		name := filepath.Join(basepath, md5[:5]+".req.txt")
		if err = os.WriteFile(name, b, 0644); err != nil {
			return nil, err
		}
		if res, err = rt.RoundTrip(req); err != nil {
			return
		}
		b, err = httputil.DumpResponse(res, true)
		if err != nil {
			return nil, err
		}
		name = filepath.Join(basepath, md5[:5]+".res.txt")
		if err = os.WriteFile(name, b, 0644); err != nil {
			return nil, err
		}
		return
	})
}
