package be

import (
	"reflect"
	"testing"
)

// DeepEqual calls t.Fatalf if want and got are different according to reflect.DeepEqual.
func DeepEqual[T any](t testing.TB, want, got T) {
	t.Helper()
	// Pass as pointers to get around the nil-interface problem
	if !reflect.DeepEqual(&want, &got) {
		t.Fatalf("reflect.DeepEqual(%#v, %#v) == false", want, got)
	}
}
