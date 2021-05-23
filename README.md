# Requests [![GoDoc](https://godoc.org/github.com/carlmjohnson/requests?status.svg)](https://godoc.org/github.com/carlmjohnson/requests) [![Go Report Card](https://goreportcard.com/badge/github.com/carlmjohnson/requests)](https://goreportcard.com/report/github.com/carlmjohnson/requests)

![Requests logo](/img/gopher-web.png)

HTTP requests for Gophers.

## Examples
```go
// Simple GET into a string
var s string
err := requests.
	URL("http://example.com").
	ToString(&s).
	Fetch(context.Background())

// Equivalent code with net/http
req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
if err != nil {
	// ...
}
res, err := http.DefaultClient.Do(req)
if err != nil {
	// ...
}
defer res.Body.Close()
b, err := io.ReadAll(res.Body)
if err != nil {
	// ...
}
s := string(b)
// 5 lines vs. 13 lines
```
```go
// Post a raw body
err := requests.
	URL("https://postman-echo.com/post").
	BodyBytes([]byte(`hello, world`)).
	ContentType("text/plain").
	Fetch(context.Background())

// Equivalent code with net/http
body := bytes.NewReader(([]byte(`hello, world`))
req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://postman-echo.com/post", body)
if err != nil {
	// ...
}
req.Header.Set("Content-Type", "text/plain")
res, err := http.DefaultClient.Do(req)
if err != nil {
	// ...
}
defer res.Body.Close()
_, err := io.ReadAll(res.Body)
if err != nil {
	// ...
}
// 5 lines vs. 14 lines
```
```go
// GET a JSON object
var post placeholder
err := requests.
	URL("https://jsonplaceholder.typicode.com").
	Path("/posts/1").
	ToJSON(&post).
	Fetch(context.Background())

// Equivalent code with net/http
var post placeholder
u, err := url.Parse("https://jsonplaceholder.typicode.com")
if err != nil {
	// ...
}
u.Path = "/posts/1"
req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), body)
if err != nil {
	// ...
}
req.Header.Set("Content-Type", "text/plain")
res, err := http.DefaultClient.Do(req)
if err != nil {
	// ...
}
defer res.Body.Close()
b, err := io.ReadAll(res.Body)
if err != nil {
	// ...
}
err := json.Unmarshal(b, &post)
if err != nil {
	// ...
}
// 6 lines vs. 23 lines
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
// net/http equivalent left as an exercise for the reader
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
```go
var params postman
err := requests.
	URL("https://postman-echo.com/get?a=1&b=2").
	Param("b", "3").
	Param("c", "4").
	Fetch(context.Background())
	// URL is https://postman-echo.com/get?a=1&b=3&c=4
```
