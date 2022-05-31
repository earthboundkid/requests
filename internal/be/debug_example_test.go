package be_test

import (
	"fmt"
	"testing"

	"github.com/carlmjohnson/requests/internal/be"
)

type mockDebug struct {
	testing.T
	failed   bool
	cleanups []func()
}

func (m *mockDebug) Run(name string, f func(t *testing.T)) {
	defer func() {
		for _, f := range m.cleanups {
			defer f()
		}
	}()
	f(&m.T)
}

func (m *mockDebug) Cleanup(f func()) {
	m.cleanups = append(m.cleanups, f)
}

func (_ *mockDebug) Log(args ...any) {
	fmt.Println(args...)
}

func (_ *mockDebug) Helper() {}

func (m *mockDebug) Fatalf(format string, args ...any) {
	m.failed = true
	fmt.Printf(format+"\n", args...)
}

func (m *mockDebug) Failed() bool { return m.failed }

func ExampleDebug() {
	t := &mockDebug{}
	// If a test fails, the callbacks will be replayed in LIFO order
	t.Run("logging-example", func(_ *testing.T) {
		x := 1
		x1 := x
		be.Debug(t, func() {
			// record some debug information about x1
			fmt.Println("x1:", x1)
		})
		x = 2
		x2 := x
		be.Debug(t, func() {
			// record some debug information about x2
			fmt.Println("x2:", x2)
		})
		be.Equal(t, x, 3)
	})
	t = &mockDebug{}
	// If a test succeeds, nothing will be replayed
	t.Run("silent-example", func(_ *testing.T) {
		y := 1
		y1 := y
		be.Debug(t, func() {
			// record some debug information about y1
			fmt.Println("y1:", y1)
		})
		y = 2
		y2 := y
		be.Debug(t, func() {
			// record some debug information about y2
			fmt.Println("y2:", y2)
		})
		be.Unequal(t, y, 3)
	})
	// Output:
	// want: 2; got: 3
	// x2: 2
	// x1: 1
}

func ExampleDebugLog() {
	t := &mockDebug{}
	// If a test fails, the logs will be replayed in LIFO order
	t.Run("logging-example", func(_ *testing.T) {
		x := 1
		be.DebugLog(t, "x: %d", x)
		x = 2
		be.DebugLog(t, "x: %d", x)
		be.Equal(t, x, 3)
	})
	t = &mockDebug{}
	// If a test succeeds, nothing will be replayed
	t.Run("silent-example", func(_ *testing.T) {
		y := 1
		be.DebugLog(t, "y: %d", y)
		y = 2
		be.DebugLog(t, "y: %d", y)
		be.Unequal(t, y, 3)
	})
	// Output:
	// want: 2; got: 3
	// x: 2
	// x: 1
}
