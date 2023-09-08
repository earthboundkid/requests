package requests

import (
	"encoding/json"
)

type jsonMarshaller = func(v any) ([]byte, error)
type jsonUnmarshaller = func(data []byte, v any) error

var (
	jsonMarshal   = json.Marshal
	jsonUnmarshal = json.Unmarshal
)

func SetJSONUnmarshaller(j jsonUnmarshaller) {
	jsonUnmarshal = j
}

func SetJSONMarshaller(j jsonMarshaller) {
	jsonMarshal = j
}

// SetJSONSerializers is a function to set global json.Marshal/json.Unmarshal function
// For faster serialization/deserialization you can use functions from, for instance, https://github.com/goccy/go-json
func SetJSONSerializers(marshaller jsonMarshaller, unmarshaller jsonUnmarshaller) {
	jsonMarshal = marshaller
	jsonUnmarshal = unmarshaller
}
