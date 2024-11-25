package requests_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
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

// TestKeepRespBodyHandlers tests the KeepRespBodyHandlers function.
func TestKeepRespBodyHandlers(t *testing.T) {
	type Common struct {
		ID int `json:"id"`
	}

	type Book struct {
		Common
		Name string `json:"name"`
	}

	var (
		book   Book
		common Common
		str    string
	)

	handler := requests.KeepRespBodyHandlers(
		requests.ToJSON(&common),
		requests.ToJSON(&book),
		requests.ToString(&str),
	)

	err := handler(&http.Response{
		Body: io.NopCloser(bytes.NewReader([]byte(`{"id":1, "name":"孙子兵法"}`))),
	})

	be.NilErr(t, err)
	be.Equal(t, 1, common.ID)
	be.Equal(t, 1, book.ID)
	be.Equal(t, "孙子兵法", book.Name)
	be.Equal(t, `{"id":1, "name":"孙子兵法"}`, str)
}
