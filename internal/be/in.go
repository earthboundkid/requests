package be

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

// In calls t.Fatalf if needle is not contained in the string or []byte haystack.
func In[byteseq ~string | ~[]byte](t testing.TB, needle string, haystack byteseq) {
	t.Helper()
	if !in(needle, haystack) {
		t.Fatalf("%q not in %q", needle, haystack)
	}
}

// NotIn calls t.Fatalf if needle is contained in the string or []byte haystack.
func NotIn[byteseq ~string | ~[]byte](t testing.TB, needle string, haystack byteseq) {
	t.Helper()
	if in(needle, haystack) {
		t.Fatalf("%q in %q", needle, haystack)
	}
}

func in[byteseq ~string | ~[]byte](needle string, haystack byteseq) bool {
	rv := reflect.ValueOf(haystack)
	switch rv.Kind() {
	case reflect.String:
		return strings.Contains(rv.String(), needle)
	case reflect.Slice:
		return bytes.Contains(rv.Bytes(), []byte(needle))
	default:
		panic("unreachable")
	}
}
