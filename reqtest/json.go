package reqtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/carlmjohnson/requests"
)

func ReplayJSON(code int, obj any) requests.Transport {
	return requests.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		data, err := json.Marshal(obj)
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
