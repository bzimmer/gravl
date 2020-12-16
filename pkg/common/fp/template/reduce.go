package template

// ReduceValueTypePtr reduces the slice to a single value.
// `f` should be a function of 2 arguments. If `val` is not supplied, returns the result of applying
// `f` to the first 2 items in `coll`, then applying `f` to that result and the 3rd item, etc. If coll
// contains no items, `f` is not called, and reduce returns either `val`, if supplied, or `nil`, if not.
// If `coll` has only 1 item, it is returned and `f` is not called.  If `val` is supplied, returns the
// result of applying `f` to `val` and the first item in `coll`, then applying `f` to that result and the
// 2nd item, etc. If `coll` contains no items, returns `val` and `f` is not called.
func ReduceValueTypePtr(f func(*ValueType, *ValueType) *ValueType, coll []*ValueType, val ...ValueType) *ValueType {
	var init *ValueType
	if len(val) > 0 {
		init = &val[0]
	}
	n := len(coll)
	if f == nil || n == 0 {
		return init
	}
	if n == 1 {
		return coll[0]
	}
	r := coll[0]
	if init != nil {
		r = f(init, r)
	}
	for i := 1; i < len(coll); i++ {
		r = f(r, coll[i])
	}
	return r
}

// ReduceValueTypePtr reduces the slice to a single value.
// `f` should be a function of 2 arguments. If `val` is not supplied, returns the result of applying
// `f` to the first 2 items in `coll`, then applying `f` to that result and the 3rd item, etc. If coll
// contains no items, `f` is not called, and reduce returns either `val`, if supplied, or `nil`, if not.
// If `coll` has only 1 item, it is returned and `f` is not called.  If `val` is supplied, returns the
// result of applying `f` to `val` and the first item in `coll`, then applying `f` to that result and the
// 2nd item, etc. If `coll` contains no items, returns `val` and `f` is not called.
func ReduceValueType(f func(ValueType, ValueType) ValueType, coll []ValueType, val ...ValueType) ValueType {
	var init ValueType
	if len(val) > 0 {
		init = val[0]
	}
	n := len(coll)
	if f == nil || n == 0 {
		return init
	}
	if n == 1 {
		return coll[0]
	}
	r := f(init, coll[0])
	for i := 1; i < len(coll); i++ {
		r = f(r, coll[i])
	}
	return r
}
