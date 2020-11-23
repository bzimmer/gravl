// Code generated by "stringer -type=Units,AlertLevel -output model_string.go"; DO NOT EDIT.

package visualcrossing

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UnitsUS-0]
	_ = x[UnitsUK-1]
	_ = x[UnitsStandard-2]
	_ = x[UnitsMetric-3]
}

const _Units_name = "UnitsUSUnitsUKUnitsStandardUnitsMetric"

var _Units_index = [...]uint8{0, 7, 14, 27, 38}

func (i Units) String() string {
	if i < 0 || i >= Units(len(_Units_index)-1) {
		return "Units(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Units_name[_Units_index[i]:_Units_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AlertLevelNone-4]
	_ = x[AlertLevelSummary-5]
	_ = x[AlertLevelDetail-6]
}

const _AlertLevel_name = "AlertLevelNoneAlertLevelSummaryAlertLevelDetail"

var _AlertLevel_index = [...]uint8{0, 14, 31, 47}

func (i AlertLevel) String() string {
	i -= 4
	if i < 0 || i >= AlertLevel(len(_AlertLevel_index)-1) {
		return "AlertLevel(" + strconv.FormatInt(int64(i+4), 10) + ")"
	}
	return _AlertLevel_name[_AlertLevel_index[i]:_AlertLevel_index[i+1]]
}