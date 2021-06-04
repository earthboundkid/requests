package requests_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/carlmjohnson/requests"
)

func TestClone(t *testing.T) {
	{
		rb1 := requests.
			URL("example.com").
			Header("a", "1").
			Header("b", "2").
			Param("a", "1").
			Param("b", "2")
		rb2 := rb1.Clone().
			Host("host.example").
			Header("b", "3").
			Header("c", "4").
			Param("b", "3").
			Param("c", "4")
		req1, err := rb1.Request(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if req1.URL.Host != "example.com" {
			t.Fatalf("bad host: %v", req1.URL)
		}
		if req1.Header.Get("b") != "2" || req1.Header.Get("c") != "" {
			t.Fatalf("bad header: %v", req1.URL)
		}
		if q := req1.URL.Query(); q.Get("b") != "2" || q.Get("c") != "" {
			t.Fatalf("bad query: %v", req1.URL)
		}
		req2, err := rb2.Request(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if req2.URL.Host != "host.example" {
			t.Fatalf("bad host: %v", req2.URL)
		}
		if req2.Header.Get("b") != "3" || req2.Header.Get("c") != "4" {
			t.Fatalf("bad header: %v", req2.URL)
		}
		if q := req2.URL.Query(); q.Get("b") != "3" || q.Get("c") != "4" {
			t.Fatalf("bad query: %v", req2.URL)
		}
	}
	{
		rb1 := new(requests.Builder).
			Host("example.com").
			Header("a", "1").
			Header("b", "2").
			Param("a", "1").
			Param("b", "2")
		rb2 := rb1.Clone().
			Host("host.example").
			Header("b", "3").
			Header("c", "4").
			Param("b", "3").
			Param("c", "4")
		req1, err := rb1.Request(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if req1.URL.Host != "example.com" {
			t.Fatalf("bad host: %v", req1.URL)
		}
		if req1.Header.Get("b") != "2" || req1.Header.Get("c") != "" {
			t.Fatalf("bad header: %v", req1.URL)
		}
		if q := req1.URL.Query(); q.Get("b") != "2" || q.Get("c") != "" {
			t.Fatalf("bad query: %v", req1.URL)
		}
		req2, err := rb2.Request(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if req2.URL.Host != "host.example" {
			t.Fatalf("bad host: %v", req2.URL)
		}
		if req2.Header.Get("b") != "3" || req2.Header.Get("c") != "4" {
			t.Fatalf("bad header: %v", req2.URL)
		}
		if q := req2.URL.Query(); q.Get("b") != "3" || q.Get("c") != "4" {
			t.Fatalf("bad query: %v", req2.URL)
		}
	}
}

func TestScheme(t *testing.T) {
	const res = `HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Date: Mon, 24 May 2021 18:48:50 GMT

An example response.`
	var s string
	const expected = `An example response.`
	var trans http.Transport
	trans.RegisterProtocol("string", requests.ReplayString(res))
	if err := requests.
		URL("example").
		Scheme("string").
		Client(&http.Client{
			Transport: &trans,
		}).
		ToString(&s).
		Fetch(context.Background()); err != nil {
		t.Fatal(err)
	}
	if s != expected {
		t.Fatalf("%q != %q", s, expected)
	}
}
