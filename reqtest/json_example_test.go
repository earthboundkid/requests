package reqtest_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/reqtest"
)

func ExampleReplayJSON() {
	// Create a ReplayJSON transport for testing
	data := map[string]any{
		"message": "Hello, World!",
		"ok":      true,
		"count":   42,
	}
	transport := reqtest.ReplayJSON(http.StatusOK, data)

	{
		type Result struct {
			Message string `json:"message"`
			OK      bool   `json:"ok"`
			Count   int    `json:"count"`
		}
		var result Result
		err := requests.
			URL("http://example.com/api").
			Transport(transport).
			ToJSON(&result).
			Fetch(context.Background())
		if err != nil {
			panic(err)
		}
		fmt.Printf("Message: %s\n", result.Message)
		fmt.Printf("Okay: %v\n", result.OK)
		fmt.Printf("Count: %d\n", result.Count)
	}

	// Modify backing map and result will change
	data["message"] = "Error!"
	data["ok"] = false
	data["count"] = 0
	{
		type Result struct {
			Message string `json:"message"`
			OK      bool   `json:"ok"`
			Count   int    `json:"count"`
		}
		var result Result
		err := requests.
			URL("http://example.com/api").
			Transport(transport).
			ToJSON(&result).
			Fetch(context.Background())
		if err != nil {
			panic(err)
		}
		fmt.Printf("Message: %s\n", result.Message)
		fmt.Printf("Okay: %v\n", result.OK)
		fmt.Printf("Count: %d\n", result.Count)
	}

	// Output:
	// Message: Hello, World!
	// Okay: true
	// Count: 42
	// Message: Error!
	// Okay: false
	// Count: 0
}
