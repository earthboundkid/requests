package requests_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/carlmjohnson/requests"
)

func BenchmarkBuilder_ToFile(b *testing.B) {
	for n := 0; n < b.N; n++ {
		d, err := os.MkdirTemp("", "to_file_*")
		if err != nil {
			b.Fatal(err)
		}
		tmpFile := filepath.Join(d, "100mb.test")
		b.Cleanup(func() {
			os.RemoveAll(d)
		})
		err = requests.URL("http://speedtest-sgp1.digitalocean.com/100mb.test").
			Client(&http.Client{Transport: http.DefaultTransport}).
			ToFile(tmpFile).
			Fetch(context.Background())
		if err != nil {
			b.Fatal(err)
		}
	}
}
