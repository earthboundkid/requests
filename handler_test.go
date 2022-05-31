package requests_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/internal/be"
)

func BenchmarkBuilder_ToFile(b *testing.B) {
	d, err := os.MkdirTemp("", "to_file_*")
	be.NilErr(b, err)
	b.Cleanup(func() {
		os.RemoveAll(d)
	})
	tmpFiles := make([]string, b.N)
	for n := 0; n < b.N; n++ {
		tmpFile := filepath.Join(d, fmt.Sprintf("10mb-%d.test", n))
		tmpFiles[n] = tmpFile
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		err = requests.URL("http://speedtest-nyc1.digitalocean.com/10mb.test").
			Client(&http.Client{Transport: http.DefaultTransport}).
			ToFile(tmpFiles[n]).
			Fetch(context.Background())
		be.NilErr(b, err)
	}
}
