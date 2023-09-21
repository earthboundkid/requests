package reqxml

import (
	"encoding/xml"

	"github.com/carlmjohnson/requests"
)

// Body is a BodyGetter that marshals a XML object.
func Body(v any) requests.BodyGetter {
	return requests.BodySerializer(xml.Marshal, v)
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
