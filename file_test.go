package requests_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/carlmjohnson/requests"
)

func TestFileTransport(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		d := t.TempDir()
		name := d + "/test.txt"
		const content = `Hello, world!`
		if err := os.WriteFile(name, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		var s string
		err := requests.
			URL(name).
			Client(&http.Client{
				Transport: requests.FileTransport,
			}).
			ToString(&s).
			Fetch(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if s != content {
			t.Fatalf("got %q", s)
		}
	})
	t.Run("404", func(t *testing.T) {
		d := t.TempDir()
		name := d + "/test.txt"
		const content = `Hello, world!`
		if err := os.WriteFile(name, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		var s string
		err := requests.
			URL(name + "x").
			Client(&http.Client{
				Transport: requests.FileTransport,
			}).
			CheckStatus(http.StatusNotFound).
			ToString(&s).
			Fetch(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(s, "no such file or directory") {
			t.Fatalf("got %q", s)
		}
	})
	t.Run("403", func(t *testing.T) {
		d := t.TempDir()
		name := d + "/test.txt"
		const content = `Hello, world!`
		if err := os.WriteFile(name, []byte(content), 0200); err != nil {
			t.Fatal(err)
		}
		var s string
		err := requests.
			URL(name).
			Client(&http.Client{
				Transport: requests.FileTransport,
			}).
			CheckStatus(http.StatusForbidden).
			ToString(&s).
			Fetch(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(s, "permission denied") {
			t.Fatalf("got %q", s)
		}
	})
	t.Run("relative", func(t *testing.T) {
		var s string
		err := requests.
			URL("img/gopher-web.acorn").
			Client(&http.Client{
				Transport: requests.FileTransport,
			}).
			AddValidator(func(res *http.Response) error {
				b, _ := io.ReadAll(res.Body)
				fmt.Println(string(b))
				return nil
			}).
			ToString(&s).
			Fetch(context.Background())
		if err != nil {
			// t.Fatal(err)
		}
		if s != "x" {
			t.Fatalf("got %q", s)
		}
	})
}
