# Requests [![GoDoc](https://godoc.org/github.com/carlmjohnson/requests?status.svg)](https://godoc.org/github.com/carlmjohnson/requests) [![Go Report Card](https://goreportcard.com/badge/github.com/carlmjohnson/requests)](https://goreportcard.com/report/github.com/carlmjohnson/requests)

HTTP requests for Gophers.

## Examples
```go
	// Simple GET into a string
	var s string
	err := requests.
		URL("http://example.com").
		ToString(&s).
		Fetch(context.Background())
```
```go
	// Post a raw body
	var data postman
	err := requests.
		URL("https://postman-echo.com/post").
		BodyBytes([]byte(`hello, world`)).
		ContentType("text/plain").
		Fetch(context.Background())
```
```go
	// GET a JSON object
	var post placeholder
	err := requests.
		URL("https://jsonplaceholder.typicode.com").
		Path("/posts/1").
		ToJSON(&post).
		Fetch(context.Background())
```
```go
	// POST a JSON object and parse the response
	var res placeholder
	req := placeholder{
		Title:  "foo",
		Body:   "baz",
		UserID: 1,
	}
	err := requests.
		URL("/posts").
		Host("jsonplaceholder.typicode.com").
		BodyJSON(&req).
		ToJSON(&res).
		Fetch(context.Background())
```
```go
	// Set headers
	var headers postman
	err := requests.
		URL("https://postman-echo.com/get").
		UserAgent("bond/james-bond").
		ContentType("secret").
		Header("martini", "shaken").
		Fetch(context.Background())
```
