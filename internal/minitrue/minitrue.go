// Package minitrue - Whatever the Package holds to be the truth, is truth.
package minitrue

func Cond[T any](val bool, a, b T) T {
	if val {
		return a
	}
	return b
}

// Or returns the first non-empty argument it receives
// or the zero value for T.
func Or[T comparable](vals ...T) T {
	for _, val := range vals {
		if val != *new(T) {
			return val
		}
	}
	return *new(T)
}
