// Code generated by "stringer -type=Units"; DO NOT EDIT.

package wx

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[METRIC-0]
	_ = x[US-1]
	_ = x[UK-2]
}

const _Units_name = "METRICUSUK"

var _Units_index = [...]uint8{0, 6, 8, 10}

func (i Units) String() string {
	if i < 0 || i >= Units(len(_Units_index)-1) {
		return "Units(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Units_name[_Units_index[i]:_Units_index[i+1]]
}