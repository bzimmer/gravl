// Code generated by "stringer -type=Format -linecomment -output=xfer_string.go"; DO NOT EDIT.

package activity

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Original-0]
	_ = x[GPX-1]
	_ = x[TCX-2]
	_ = x[FIT-3]
}

const _Format_name = "originalgpxtcxfit"

var _Format_index = [...]uint8{0, 8, 11, 14, 17}

func (i Format) String() string {
	if i < 0 || i >= Format(len(_Format_index)-1) {
		return "Format(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Format_name[_Format_index[i]:_Format_index[i+1]]
}