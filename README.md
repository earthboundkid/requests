# Requests [![GoDoc](https://godoc.org/github.com/carlmjohnson/requests?status.svg)](https://godoc.org/github.com/carlmjohnson/requests) [![Go Report Card](https://goreportcard.com/badge/github.com/carlmjohnson/requests)](https://goreportcard.com/report/github.com/carlmjohnson/requests)

![Requests logo](/img/gopher-web.png)

## _HTTP requests for Gophers._

**The problem**: Go's net/http is powerful and versatile, but using it correctly for client requests can be extremely verbose.

**The solution**: The requests.Builder type is a convenient way to build, send, and handle HTTP requests. Builder has a fluent API with methods returning a pointer to the same struct, which allows for declaratively describing a request by method chaining.

Requests also comes with tools for building custom http transports, include a request recorder and replayer for testing.

## Examples
### Simple GET into a string

```go
// code with net/http
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

// equivalent code using requests
var s string
err := requests.
	URL("http://example.com").
	ToString(&s).
	Fetch(context.Background())

// 5 lines vs. 13 lines
```

### POST a raw body

```go
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

### GET a JSON object

```go
var post placeholder
err := requests.
	URL("https://jsonplaceholder.typicode.com").
	Pathf("/posts/%d", 1).
	ToJSON(&post).
	Fetch(context.Background())

// Equivalent code with net/http
var post placeholder
u, err := url.Parse("https://jsonplaceholder.typicode.com")
if err != nil {
	// ...
}
u.Path = fmt.Sprintf("/posts/%d", 1)
req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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
err := json.Unmarshal(b, &post)
if err != nil {
	// ...
}
// 6 lines vs. 23 lines
```

### POST a JSON object and parse the response

```go
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

### Set custom headers for a request

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

### Easily manipulate query parameters

```go
var params postman
err := requests.
	URL("https://postman-echo.com/get?a=1&b=2").
	Param("b", "3").
	Param("c", "4").
	Fetch(context.Background())
	// URL is https://postman-echo.com/get?a=1&b=3&c=4
```

### Record and replay responses

```go
// record a request to the file system
cl.Transport = requests.Record(nil, "somedir")
var s1, s2 string
err := requests.URL("http://example.com").
	Client(&cl).
	ToString(&s1).
	Fetch(context.Background())
check(err)

// now replay the request in tests
cl.Transport = requests.Replay("somedir")
err = requests.URL("http://example.com").
	Client(&cl).
	ToString(&s2).
	Fetch(context.Background())
check(err)
assert(s1 == s2) // true
```
