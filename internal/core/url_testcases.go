package core

// PathCases is exported to share tests with requests
var PathCases = map[string]struct {
	Base   string
	Paths  []string
	Result string
}{
	"base-only": {
		"example",
		[]string{},
		"https://example",
	},
	"base+abspath": {
		"https://example",
		[]string{"/a"},
		"https://example/a",
	},
	"multi-abs-paths": {
		"https://example",
		[]string{"/a", "/b/", "/c"},
		"https://example/c",
	},
	"base+rel-path": {
		"https://example/a/",
		[]string{"./b"},
		"https://example/a/b",
	},
	"base+rel-paths": {
		"https://example/a/",
		[]string{"./b/", "./c"},
		"https://example/a/b/c",
	},
	"rel-path": {
		"https://example/",
		[]string{"a/", "./b"},
		"https://example/a/b",
	},
	"base+multi-paths": {
		"https://example/a/",
		[]string{"b/", "c"},
		"https://example/a/b/c",
	},
	"base+slash+multi-paths": {
		"https://example/a/",
		[]string{"b/", "c"},
		"https://example/a/b/c",
	},
	"multi-root": {
		"https://example/",
		[]string{"a", "b", "c"},
		"https://example/c",
	},
	"dot-dot-paths": {
		"https://example/",
		[]string{"a/", "b/", "../c"},
		"https://example/a/c",
	},
	"more-dot-dot-paths": {
		"https://example/",
		[]string{"a/b/c/", "../d/", "../e"},
		"https://example/a/b/e",
	},
	"more-dot-dot-paths+rel-path": {
		"https://example/",
		[]string{"a/b/c/", "../d/", "../e/", "./f"},
		"https://example/a/b/e/f",
	},
	"even-more-dot-dot-paths+base": {
		"https://example/a/b/c/",
		[]string{"../../d"},
		"https://example/a/d",
	},
	"too-many-dot-dot-paths": {
		"https://example",
		[]string{"../a"},
		"https://example/a",
	},
	"too-many-dot-dot-paths+base": {
		"https://example/",
		[]string{"../a"},
		"https://example/a",
	},
	"last-abs-path-wins": {
		"https://example/a/",
		[]string{"b/", "c/", "/d"},
		"https://example/d",
	},
}
