package requests

import (
	"bytes"
	"io"
	"net/url"
	"os"
	"strings"
)

// BodyGetter provides a Builder with a source for a request body.
type BodyGetter = func() (io.ReadCloser, error)

// bodyReader is a BodyGetter that returns an io.Reader.
func bodyReader(r io.Reader) BodyGetter {
	return func() (io.ReadCloser, error) {
		if rc, ok := r.(io.ReadCloser); ok {
			return rc, nil
		}
		return rc(r), nil
	}
}

// bodyWriter is a BodyGetter that pipes writes into a request body.
func bodyWriter(f func(w io.Writer) error) BodyGetter {
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

// bodyBytes is a BodyGetter that returns the provided raw bytes.
func bodyBytes(b []byte) BodyGetter {
	return func() (io.ReadCloser, error) {
		return rc(bytes.NewReader(b)), nil
	}
}

// BodySerializer is a BodyGetter
// that uses the provided [Serializer]
// to build the body of a request from v.
func BodySerializer(s Serializer, v any) BodyGetter {
	return func() (io.ReadCloser, error) {
		b, err := s(v)
		if err != nil {
			return nil, err
		}
		return rc(bytes.NewReader(b)), nil
	}
}

// bodyJSON is a [BodySerializer]
// that uses [JSONSerializer] to marshal the object.
func bodyJSON(v any) BodyGetter {
	return BodySerializer(JSONSerializer, v)
}

// bodyForm is a BodyGetter that builds an encoded form body.
func bodyForm(data url.Values) BodyGetter {
	return func() (r io.ReadCloser, err error) {
		return rc(strings.NewReader(data.Encode())), nil
	}
}

// bodyFile is a BodyGetter that reads the provided file path.
func bodyFile(name string) BodyGetter {
	return func() (r io.ReadCloser, err error) {
		return os.Open(name)
	}
}
