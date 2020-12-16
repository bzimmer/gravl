package template

func EveryValueTypePtr(f func(*ValueType) bool, coll []*ValueType) bool {
	if f == nil || len(coll) == 0 {
		return false
	}
	for _, v := range coll {
		if !f(v) {
			return false
		}
	}
	return true
}

func EveryValueType(f func(ValueType) bool, coll []ValueType) bool {
	if f == nil || len(coll) == 0 {
		return false
	}
	for _, v := range coll {
		if !f(v) {
			return false
		}
	}
	return true
}
