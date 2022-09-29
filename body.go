package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"strings"
)

// nopCloser is like io.NopCloser(),
// but it is a concrete type so we can strip it out
// before setting a body on a request.
// See https://github.com/carlmjohnson/requests/discussions/49
type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

var _ io.ReadCloser = nopCloser{}

// BodyGetter provides a Builder with a source for a request body.
type BodyGetter = func() (io.ReadCloser, error)

// BodyReader is a BodyGetter that returns an io.Reader.
func BodyReader(r io.Reader) BodyGetter {
	return func() (io.ReadCloser, error) {
		if rc, ok := r.(io.ReadCloser); ok {
			return rc, nil
		}
		return nopCloser{r}, nil
	}
}

// BodyWriter is a BodyGetter that pipes writes into a request body.
func BodyWriter(f func(w io.Writer) error) BodyGetter {
	return func() (io.ReadCloser, error) {
		r, w := io.Pipe()
		go func() {
			var err error
			defer func() {
				w.CloseWithError(err)
			}()
			err = f(w)
		}()
		return r, nil
	}
}

// BodyBytes is a BodyGetter that returns the provided raw bytes.
func BodyBytes(b []byte) BodyGetter {
	return func() (io.ReadCloser, error) {
		return nopCloser{bytes.NewReader(b)}, nil
	}
}

// BodyJSON is a BodyGetter that marshals a JSON object.
func BodyJSON(v any) BodyGetter {
	return func() (io.ReadCloser, error) {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return nopCloser{bytes.NewReader(b)}, nil
	}
}

// BodyForm is a BodyGetter that builds an encoded form body.
func BodyForm(data url.Values) BodyGetter {
	return func() (r io.ReadCloser, err error) {
		return nopCloser{strings.NewReader(data.Encode())}, nil
	}
}

// BodyFile is a BodyGetter that reads the provided file path.
func BodyFile(name string) BodyGetter {
	return func() (r io.ReadCloser, err error) {
		return os.Open(name)
	}
}
