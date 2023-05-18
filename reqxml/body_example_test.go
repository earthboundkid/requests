package reqxml_test

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http/httputil"
	"strings"

	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/reqxml"
)

func ExampleBody() {
	// Make an XML object
	type Object struct {
		XMLName xml.Name `xml:"object"`
		ID      int      `xml:"id,attr"`
		Verbose bool     `xml:"communication>mode>verbosity>high"`
		Comment string   `xml:",comment"`
	}

	v := &Object{
		ID:      42,
		Verbose: true,
		Comment: "Hello!",
	}

	req, err := requests.
		URL("http://example.com").
		// reqxml.BodyConfig sets the content type for us
		Config(reqxml.BodyConfig(&v)).
		Request(context.Background())
	if err != nil {
		fmt.Println("Error!", err)
	}

	// Pretty print the request
	b, _ := httputil.DumpRequestOut(req, true)
	fmt.Println(strings.ReplaceAll(string(b), "\r", ""))

	// Output:
	// POST / HTTP/1.1
	// Host: example.com
	// User-Agent: Go-http-client/1.1
	// Content-Length: 122
	// Content-Type: application/xml
	// Accept-Encoding: gzip
	//
	// <object id="42"><communication><mode><verbosity><high>true</high></verbosity></mode></communication><!--Hello!--></object>
}
