package requests_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/gophertools/requests"
	"github.com/gophertools/requests/internal/be"
	"github.com/gophertools/requests/reqtest"
)

func TestJSONV1(t *testing.T) {
	t.Parallel()
	type obj struct {
		Name  string
		Count int `json:"count,omitempty"`
	}
	ctx := context.Background()

	t.Run("BodyJSON", func(t *testing.T) {
		req, err := requests.New().BodyJSON(obj{Name: "a"}).Request(ctx)
		be.NilErr(t, err)
		b, err := io.ReadAll(req.Body)
		be.NilErr(t, err)
		be.Equal(t, `{"Name":"a","count":0}`, string(b))
	})
	t.Run("BodyJSONV1", func(t *testing.T) {
		req, err := requests.New().BodyJSONV1(obj{Name: "a"}).Request(ctx)
		be.NilErr(t, err)
		be.Equal(t, "application/json", req.Header.Get("Content-Type"))
		b, err := io.ReadAll(req.Body)
		be.NilErr(t, err)
		be.Equal(t, `{"Name":"a"}`, string(b))
	})
	t.Run("ToJSON", func(t *testing.T) {
		var got obj
		err := requests.New().
			Transport(reqtest.ReplayString("HTTP/1.1 200 OK\n\n{\"name\": \"a\"}")).
			ToJSON(&got).
			Fetch(ctx)
		be.NilErr(t, err)
		be.Equal(t, "", got.Name)
	})
	t.Run("ToJSONV1", func(t *testing.T) {
		var got obj
		err := requests.New().
			Transport(reqtest.ReplayString("HTTP/1.1 200 OK\n\n{\"name\": \"a\"}")).
			ToJSONV1(&got).
			Fetch(ctx)
		be.NilErr(t, err)
		be.Equal(t, "a", got.Name)
	})
	t.Run("ErrorJSONV1", func(t *testing.T) {
		var good, bad obj
		err := requests.New().
			Transport(reqtest.ReplayString("HTTP/1.1 418 I'm a teapot\n\n{\"name\": \"brewing\"}")).
			ToJSONV1(&good).
			ErrorJSONV1(&bad).
			Fetch(ctx)
		be.True(t, errors.Is(err, requests.ErrInvalidHandled))
		be.Equal(t, "brewing", bad.Name)
		be.Equal(t, "", good.Name)
	})
}
