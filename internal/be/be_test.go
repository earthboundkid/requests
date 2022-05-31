package be_test

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/carlmjohnson/requests/internal/be"
)

type testingTB struct {
	testing.TB
	failed bool
	w      io.Writer
}

func (t *testingTB) Helper() {}

func (t *testingTB) Fatalf(format string, args ...any) {
	t.failed = true
	fmt.Fprintf(t.w, format, args...)
}

func Test(t *testing.T) {
	beOkay := func(callback func(tb testing.TB)) {
		t.Helper()
		var buf strings.Builder
		tb := &testingTB{w: &buf}
		callback(tb)
		if tb.failed {
			t.Fatal("failed too soon")
		}
		if buf.String() != "" {
			t.Fatal("wrote too much")
		}
	}
	beOkay(func(tb testing.TB) { be.Zero(tb, time.Time{}.Local()) })
	beOkay(func(tb testing.TB) { be.Zero(tb, []string(nil)) })
	beOkay(func(tb testing.TB) { be.Nonzero(tb, []string{""}) })
	beOkay(func(tb testing.TB) { be.NilErr(tb, nil) })
	beOkay(func(tb testing.TB) { be.True(tb, true) })
	beOkay(func(tb testing.TB) { be.False(tb, false) })
	beBad := func(callback func(tb testing.TB)) {
		t.Helper()
		var buf strings.Builder
		tb := &testingTB{w: &buf}
		callback(tb)
		if !tb.failed {
			t.Fatal("did not fail")
		}
		if buf.String() == "" {
			t.Fatal("wrote too little")
		}
	}
	beBad(func(tb testing.TB) { be.AllEqual(tb, []string{}, []string{""}) })
	beBad(func(tb testing.TB) { be.Nonzero(tb, time.Time{}.Local()) })
	beBad(func(tb testing.TB) { be.Zero(tb, []string{""}) })
	beBad(func(tb testing.TB) { be.Nonzero(tb, []string(nil)) })
	beBad(func(tb testing.TB) { be.NilErr(tb, errors.New("")) })
	beBad(func(tb testing.TB) { be.True(tb, false) })
	beBad(func(tb testing.TB) { be.False(tb, true) })
}
