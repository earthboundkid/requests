package reqxml

import (
	"encoding/xml"

	"github.com/gophertools/requests"
)

// To decodes a response as an XML object.
func To(v any) requests.ResponseHandler {
	return requests.ToDeserializer(xml.Unmarshal, v)
}
