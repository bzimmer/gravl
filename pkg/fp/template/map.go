package template

// MapValueTypePtr iterates over the slice of ValueType returning a new slice after applying the function
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

// func MapValueType(f func(ValueType) ValueType, coll []ValueType) []ValueType {
// 	if f == nil {
// 		return []ValueType{}
// 	}
// 	w := make([]ValueType, len(coll))
// 	for i, v := range coll {
// 		w[i] = f(v)
// 	}
// 	return w
// }
