package reqtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/carlmjohnson/requests"
)

// ReplayJSON returns a [requests.Transport]
// that always responds with the given object marshaled as JSON
// and the provided HTTP status code.
// The object is marshaled on each request,
// so modifications to mutable objects will be reflected in subsequent responses.
//
// If the object cannot be marshaled to JSON, the transport returns the
// wrapped marshaling error.
func ReplayJSON(code int, obj any) requests.Transport {
	return requests.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		data, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal error in ReplayJSON: %v", err)
		}

		return &http.Response{
			Status:     fmt.Sprintf("%d %s", code, http.StatusText(code)),
			StatusCode: code,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body:          io.NopCloser(bytes.NewReader(data)),
			ContentLength: int64(len(data)),
			Request:       req,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
		}, nil
	})
}
