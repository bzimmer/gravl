package template

func FilterValueType(f func(ValueType) bool, coll []ValueType) []ValueType {
	if f == nil {
		return []ValueType{}
	}
	var w []ValueType
	for _, v := range coll {
		if f(v) {
			w = append(w, v)
		}
	}
	return w
}

func FilterValueTypePtr(f func(*ValueType) bool, coll []*ValueType) []*ValueType {
	if f == nil {
		return []*ValueType{}
	}
	var w []*ValueType
	for _, v := range coll {
		if f(v) {
			w = append(w, v)
		}
	}
	return w
}
