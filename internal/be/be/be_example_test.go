package be_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/carlmjohnson/requests/internal/be"
)

type mockingT struct{ *testing.T }

func (_ mockingT) Helper() {}

func (_ mockingT) Fatalf(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}

func Example() {
	// mock *testing.T for example purposes
	var t mockingT

	be.Equal(t, "hello", "world")     // bad
	be.Equal(t, "goodbye", "goodbye") // good

	be.Unequal(t, "hello", "world")     // good
	be.Unequal(t, "goodbye", "goodbye") // bad

	s := []int{1, 2, 3}
	be.AllEqual(t, []int{1, 2, 3}, s) // good
	be.AllEqual(t, []int{3, 2, 1}, s) // bad

	var err error
	be.NilErr(t, err)  // good
	be.Nonzero(t, err) // bad
	err = errors.New("(O_o)")
	be.NilErr(t, err)  // bad
	be.Nonzero(t, err) // good

	type mytype string
	var mystring mytype = "hello, world"
	be.In(t, "world", mystring)                 // good
	be.In(t, "World", mystring)                 // bad
	be.NotIn(t, "\x01", []byte("\a\b\x00\r\t")) // good
	be.NotIn(t, "\x00", []byte("\a\b\x00\r\t")) // bad

	// Output:
	// want: hello; got: world
	// got: goodbye
	// want: [3 2 1]; got: [1 2 3]
	// got: <nil>
	// got: (O_o)
	// "World" not in "hello, world"
	// "\x00" in "\a\b\x00\r\t"
}
