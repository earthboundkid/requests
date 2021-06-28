package requests

import (
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// RegisterFileProtocol registers FileTransport as the handler for URLs with
// "file" or "" scheme. If t is nil, it attempts to use http.DefaultTransport,
// which may panic if it it is not a *http.Transport. This may be a security
// risk in some environments. Use responsibly.
func RegisterFileProtocol(t *http.Transport) {
	if t == nil {
		t = http.DefaultTransport.(*http.Transport)
	}
	t.RegisterProtocol("", RoundTripFunc(nil))
	t.RegisterProtocol("file", RoundTripFunc(nil))
}

// FileTransport returns local files as responses. It is like
// http.NewFileTransport but uses the current working directory for relative
// paths and / for absolute paths instead of being limited to an
// http.FileSystem.
var FileTransport http.RoundTripper = RoundTripFunc(fileTransport)

func fileTransport(req *http.Request) (res *http.Response, err error) {
	res = &http.Response{
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		Header:        http.Header{},
		Body:          nil,
		ContentLength: 0,
		Request:       req,
	}
	log.Printf("%v", req.URL)
	f, err := os.Open(req.URL.Path)
	if err != nil {
		res.Body = io.NopCloser(strings.NewReader(err.Error()))
		res.StatusCode = http.StatusInternalServerError
		if os.IsNotExist(err) {
			res.StatusCode = http.StatusNotFound

			err = nil
		} else if os.IsPermission(err) {
			res.StatusCode = http.StatusForbidden
			err = nil
		}
		res.Status = http.StatusText(res.StatusCode)
		return
	}
	res.StatusCode = http.StatusOK
	res.Status = http.StatusText(res.StatusCode)
	res.Body = f
	if stat, err2 := f.Stat(); err2 == nil {
		res.ContentLength = stat.Size()
	}
	ctype := mime.TypeByExtension(filepath.Ext(req.URL.Path))
	if ctype == "" {
		// read a chunk to decide between utf-8 text and binary
		var buf [512]byte
		n, _ := io.ReadFull(f, buf[:])
		ctype = http.DetectContentType(buf[:n])
		_, _ = f.Seek(0, io.SeekStart) // rewind to output whole file
	}
	res.Header.Set("Content-Type", ctype)
	return
}
