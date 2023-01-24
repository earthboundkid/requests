package requests

import (
	"errors"
	"fmt"
)

type multierror struct {
	msg  string
	errs []error
}

func joinerrs(a, b error, format string, args ...any) error {
	switch {
	case a == nil && b == nil:
		return nil
	case a != nil && b == nil:
		return a
	case a == nil && b != nil:
		return b
	case a != nil && b != nil:
		return multierror{
			msg:  fmt.Sprintf(format, args...),
			errs: []error{a, b},
		}
	}
	panic("unreachable")
}

func (m multierror) Is(target error) bool {
	for _, err := range m.errs {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (m multierror) As(target any) bool {
	for _, err := range m.errs {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

func (m multierror) Error() string {
	return m.msg
}
