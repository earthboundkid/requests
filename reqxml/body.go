package reqxml

import (
	"encoding/xml"
	"io"

	"github.com/carlmjohnson/requests"
)

// Body is a BodyGetter that marshals a XML object.
func Body(v any) requests.BodyGetter {
	return func() (io.ReadCloser, error) {
		b, err := xml.Marshal(v)
		if err != nil {
			return nil, err
		}
		return requests.BodyBytes(b)()
	}
}

// BodyConfig sets the Builder's request body to the marshaled XML.
// It also sets ContentType to "application/xml".
func BodyConfig(v any) requests.Config {
	return func(rb *requests.Builder) {
		rb.
			Body(Body(v)).
			ContentType("application/xml")
	}
}
