package core_test

import (
	"testing"

	"github.com/carlmjohnson/requests/internal/be"
	"github.com/carlmjohnson/requests/internal/core"
)

func TestPath(t *testing.T) {
	t.Parallel()
	for name, tc := range core.PathCases {
		t.Run(name, func(t *testing.T) {
			var b core.URLBuilder
			b.BaseURL(tc.Base)
			for _, p := range tc.Paths {
				b.Path(p)
			}
			u, err := b.URL()
			be.NilErr(t, err)
			be.Equal(t, tc.Result, u.String())
		})
	}
}
