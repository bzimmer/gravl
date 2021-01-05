package template

func MapValueType(f func(ValueType) ValueType, coll []ValueType) []ValueType {
	if f == nil {
		return []ValueType{}
	}
	w := make([]ValueType, len(coll))
	for i, v := range coll {
		w[i] = f(v)
	}
	return w
}

func MapValueTypePtr(f func(*ValueType) *ValueType, coll []*ValueType) []*ValueType {
	if f == nil {
		return []*ValueType{}
	}
	w := make([]*ValueType, len(coll))
	for i, v := range coll {
		w[i] = f(v)
	}
	return w
}
