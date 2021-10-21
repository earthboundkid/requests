package requests_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/carlmjohnson/requests"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ExampleToBufioReader() {
	// read a response line by line for a sentinel
	found := false
	err := requests.
		URL("http://example.com").
		Handle(requests.ToBufioReader(func(r *bufio.Reader) error {
			var err error
			for s := ""; err == nil; {
				if strings.Contains(s, "Example Domain") {
					found = true
					return nil
				}
				// read one line from response
				s, err = r.ReadString('\n')
			}
			if err == io.EOF {
				return nil
			}
			return err
		})).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to example.com:", err)
	}
	fmt.Println(found)
	// Output:
	// true
}

func ExampleToBufioScanner() {
	// read a response line by line for a sentinel
	found := false
	needle := []byte("Example Domain")
	err := requests.
		URL("http://example.com").
		Handle(requests.ToBufioScanner(func(s *bufio.Scanner) error {
			// read one line at time from response
			for s.Scan() {
				if bytes.Contains(s.Bytes(), needle) {
					found = true
					return nil
				}
			}
			return s.Err()
		})).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to example.com:", err)
	}
	fmt.Println(found)
	// Output:
	// true
}

func ExampleToHTML() {
	var doc html.Node
	err := requests.
		URL("http://example.com").
		Handle(requests.ToHTML(&doc)).
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
