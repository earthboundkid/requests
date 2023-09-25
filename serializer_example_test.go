package requests_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/carlmjohnson/requests"
)

func ExampleJSONSerializer() {
	data := struct {
		A string `json:"a"`
		B int    `json:"b"`
		C []bool `json:"c"`
	}{
		"Hello", 42, []bool{true, false},
	}

	// A request using the default serializer
	req, err := requests.
		New().
		BodyJSON(&data).
		Request(context.Background())
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, req.Body)
	fmt.Println()

	// Restore default after this test
	defaultSerializer := requests.JSONSerializer
	defer func() {
		requests.JSONSerializer = defaultSerializer
	}()

	// Serialize with indented JSON
	requests.JSONSerializer = func(v any) ([]byte, error) {
		return json.MarshalIndent(v, "", "  ")
	}
	req, err = requests.
		New().
		BodyJSON(&data).
		Request(context.Background())
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, req.Body)

	// Output:
	// {"a":"Hello","b":42,"c":[true,false]}
	// {
	//   "a": "Hello",
	//   "b": 42,
	//   "c": [
	//     true,
	//     false
	//   ]
	// }
}
