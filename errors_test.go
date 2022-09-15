package requests_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"testing/fstest"

	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/internal/be"
)

func TestOnError(t *testing.T) {
	fsys := fstest.MapFS{
		"fsys.example - MKIYDwjs.res.txt": &fstest.MapFile{
			Data: []byte(`HTTP/1.1 400 BadRequest
Content-Type: application/json; charset=UTF-8
Date: Mon, 24 May 2021 18:48:50 GMT

{"msg": "you messed up"}`),
		},
	}

	const (
		defaultErrorStr = "response error for http://fsys.example: unexpected status: 400"
		url = "http://fsys.example"
	)

	t.Run("no_error_handler_returns_default", func(t *testing.T) {
		err := requests.
			URL(url).
			Client(&http.Client{
				Transport: requests.ReplayFS(fsys),
			}).
			Fetch(context.Background())

		be.Equal(t, defaultErrorStr, err.Error())
	})

	t.Run("set_validation_handler", func(t *testing.T) {
		didCall := false
		err := requests.
			URL(url).
			Client(&http.Client{
				Transport: requests.ReplayFS(fsys),
			}).
			OnError(requests.ErrorKindValidationErr, func(kind requests.ErrorKind, cause error, resp *http.Response) error {
				didCall = true
				return cause
			}).
			Fetch(context.Background())

		be.Equal(t, true, didCall)
		be.Equal(t, defaultErrorStr, err.Error())
	})

	t.Run("set_validation_handler", func(t *testing.T) {
		didCall := false
		err := requests.
			URL(url).
			Client(&http.Client{
				Transport: requests.ReplayFS(fsys),
			}).
			AddValidator(nil).
			Handle(func(r *http.Response) error {
				return fmt.Errorf("some-handling-error")
			}).
			OnError(requests.ErrorKindHandlerErr, func(kind requests.ErrorKind, cause error, resp *http.Response) error {
				didCall = true
				return cause
			}).
			Fetch(context.Background())

		be.Equal(t, true, didCall)
		be.Equal(t, "some-handling-error", err.Error())
	})
}
