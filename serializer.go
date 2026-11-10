package requests

import (
	jsonv1 "encoding/json"
	jsonv2 "encoding/json/v2"
)

// Serializer is a function
// that can convert arbitrary data
// to bytes in some format.
type Serializer = func(v any) ([]byte, error)

// Deserializer is a function
// that can read data in some format
// and store the result in v.
type Deserializer = func(data []byte, v any) error

var (
	// JSONV1Serializer is [jsonv1.Marshal].
	JSONV1Serializer Serializer = jsonv1.Marshal
	// JSONSerializer is used by [BodyJSON] and [Builder.BodyJSON].
	JSONSerializer Serializer = func(v any) ([]byte, error) {
		return jsonv2.Marshal(v)
	}
	// JSONDeserializer is used by [ToJSON] and [Builder.ToJSON].
	// The default deserializer may be changed in a future version of requests.
	JSONDeserializer Deserializer = func(data []byte, v any) error {
		return jsonv2.Unmarshal(data, v)
	}
)
