package be

import (
	"fmt"
	"testing"
)

// Debug takes a callback that will only be run after the test fails.
func Debug(t testing.TB, f func()) {
	t.Helper()
	t.Cleanup(func() {
		t.Helper()
		if t.Failed() {
			f()
		}
	})
}

// DebugLog records a message that will only be logged after the test fails.
func DebugLog(t testing.TB, format string, args ...any) {
	t.Helper()
	msg := fmt.Sprintf(format, args...)
	t.Cleanup(func() {
		t.Helper()
		if t.Failed() {
			t.Log(msg)
		}
	})
}
