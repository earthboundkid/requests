package util

func Clip[T any](sp *[]T) {
	s := *sp
	*sp = s[:len(s):len(s)]
}

func Cond[T any](val bool, a, b T) T {
	if val {
		return a
	}
	return b
}

func First[T comparable](a, b T) T {
	return Cond(a != *new(T), a, b)
}
