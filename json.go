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

func SetJSONSerializers(marshaller jsonMarshaller, unmarshaller jsonUnmarshaller) {
	jsonMarshal = marshaller
	jsonUnmarshal = unmarshaller
}
