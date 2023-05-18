package reqxml

import (
	"encoding/xml"
	"io"
	"net/http"

	"github.com/carlmjohnson/requests"
)

// To decodes a response as an XML object.
func To(v any) requests.ResponseHandler {
	return func(res *http.Response) error {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if err = xml.Unmarshal(data, v); err != nil {
			return err
		}
		return nil
	}
}
