package requests

import (
	"compress/gzip"
	"io"
)

// Config allows Builder to be extended by setting several options at once.
// For example, a Config might set a Body and its ContentType.
type Config = func(rb *Builder)

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
