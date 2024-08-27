package reqhtml_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/reqhtml"
	"github.com/carlmjohnson/requests/reqtest"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func init() {
	http.DefaultTransport = reqtest.ReplayString(`HTTP/1.1 200 OK

	<a href="https://www.iana.org/domains/example"></a>`)
}

func ExampleTo() {
	var doc html.Node
	err := requests.
		URL("http://example.com").
		Handle(reqhtml.To(&doc)).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to example.com:", err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.DataAtom == atom.A {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					fmt.Println("link:", attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(&doc)
	// Output:
	// link: https://www.iana.org/domains/example
}

func ExampleBody() {
	link := html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			{Key: "href", Val: "http://example.com"},
		},
	}
	text := html.Node{
		Type: html.TextNode,
		Data: "Hello, World!",
	}
	link.AppendChild(&text)

	req, err := requests.
		URL("http://example.com").
		Config(reqhtml.Body(&link)).
		Request(context.Background())
	if err != nil {
		panic(err)
	}
	b, err := httputil.DumpRequest(req, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%q\n", b)

	// Output:
	// "POST / HTTP/1.1\r\nHost: example.com\r\nContext-Type: text/html\r\n\r\n<a href=\"http://example.com\">Hello, World!</a>"
}
