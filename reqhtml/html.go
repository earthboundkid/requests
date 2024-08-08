// Package reqhtml contains utilities for sending and receiving x/net/html objects.
package reqhtml

import (
	"io"

	"github.com/carlmjohnson/requests"
	"golang.org/x/net/html"
)

// To decodes a response as an html document.
func To(n *html.Node) requests.ResponseHandler {
	return requests.ToHTML(n)
}

// Body sets the requests.Builder's request body to the HTML document.
// It also sets ContentType to "text/html"
// if it is not otherwise set.
func Body(n *html.Node) requests.Config {
	return func(rb *requests.Builder) {
		rb.
			Body(requests.BodyWriter(func(w io.Writer) error {
				return html.Render(w, n)
			})).
			HeaderOptional("context-type", "text/html")
	}
}
