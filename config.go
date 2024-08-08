package requests

import (
	"compress/gzip"
	"io"
	"net/http/httptest"
)

// Config allows Builder to be extended by setting several options at once.
// For example, a Config might set a Body and its ContentType.
type Config = func(rb *Builder)

// GzipConfig writes a gzip stream to its request body using a callback.
// It also sets the appropriate Content-Encoding header and automatically
// closes and the stream when the callback returns.
func GzipConfig(level int, h func(gw *gzip.Writer) error) Config {
	return func(rb *Builder) {
		rb.
			Header("Content-Encoding", "gzip").
			BodyWriter(func(w io.Writer) error {
				gw, err := gzip.NewWriterLevel(w, level)
				if err != nil {
					return err
				}
				if err = h(gw); err != nil {
					gw.Close()
					return err
				}
				return gw.Close()
			})
	}
}

// TestServerConfig returns a Config
// which sets the Builder's BaseURL to s.URL
// and the Builder's Client to s.Client().
//
// Deprecated: Use reqtest.Server.
func TestServerConfig(s *httptest.Server) Config {
	return func(rb *Builder) {
		rb.
			BaseURL(s.URL).
			Client(s.Client())
	}
}
