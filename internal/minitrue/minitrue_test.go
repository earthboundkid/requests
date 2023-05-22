package minitrue_test

import (
	"testing"

	"github.com/carlmjohnson/requests/internal/be"
	"github.com/carlmjohnson/requests/internal/minitrue"
)

func TestCond(t *testing.T) {
	be.Equal(t, 1, minitrue.Cond(true, 1, 2))
	be.Equal(t, 2, minitrue.Cond(false, 1, 2))
}

func TestOr(t *testing.T) {
	be.Equal(t, 0, minitrue.Or[int]())
	be.Equal(t, 0, minitrue.Or(0))
	be.Equal(t, 1, minitrue.Or(1))
	be.Equal(t, 2, minitrue.Or(0, 2))
	be.Equal(t, 3, minitrue.Or(0, 0, 3))
}
