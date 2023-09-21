package requests

import (
	"encoding/json"
)

type Serializer = func(v any) ([]byte, error)
type Deserializer = func(data []byte, v any) error

var (
	JSONSerializer   = json.Marshal
	JSONDeserializer = json.Unmarshal
)
