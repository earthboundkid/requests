package be_test

import (
	"github.com/carlmjohnson/requests/internal/be"
)

func ExampleDeepEqual() {
	var t mockingT

	// good
	m1 := map[int]bool{1: true, 2: false}
	m2 := map[int]bool{1: true, 2: false}
	be.DeepEqual(t, m1, m2)

	// bad
	var s1 []int
	s2 := []int{}
	be.DeepEqual(t, s1, s2) // DeepEqual is picky about nil vs. len 0

	// Output:
	// reflect.DeepEqual([]int(nil), []int{}) == false
}
