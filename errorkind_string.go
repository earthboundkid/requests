// Code generated by "stringer -type=ErrorKind"; DO NOT EDIT.

package requests

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ErrorKindNone-0]
	_ = x[ErrorKindUnknown-1]
	_ = x[ErrorKindURLParse-2]
	_ = x[ErrorKindBodyGet-3]
	_ = x[ErrorKindUnknownMethod-4]
	_ = x[ErrorKindNilContext-5]
	_ = x[ErrorKindConnection-6]
	_ = x[ErrorKindValidator-7]
	_ = x[ErrorKindHandler-8]
}

const _ErrorKind_name = "ErrorKindNoneErrorKindUnknownErrorKindURLParseErrorKindBodyGetErrorKindUnknownMethodErrorKindNilContextErrorKindConnectionErrorKindValidatorErrorKindHandler"

var _ErrorKind_index = [...]uint8{0, 13, 29, 46, 62, 84, 103, 122, 140, 156}

func (i ErrorKind) String() string {
	if i < 0 || i >= ErrorKind(len(_ErrorKind_index)-1) {
		return "ErrorKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ErrorKind_name[_ErrorKind_index[i]:_ErrorKind_index[i+1]]
}
