package template

// GroupByKeyTypeValueTypePtr returns a mapping of KeyType to []*ValueType based on the return
// value of function `f`
//
// Example:
//   groups := GroupByKeyTypeValueTypePtr(func(val *ValueType) KeyType {
// 	   return KeyTypeIsh
//   }, vals)
func GroupByKeyTypeValueTypePtr(f func(act *ValueType) KeyType, coll []*ValueType) map[KeyType][]*ValueType {
	res := make(map[KeyType][]*ValueType)
	EveryValueTypePtr(func(act *ValueType) bool {
		key := f(act)
		res[key] = append(res[key], act)
		return true
	}, coll)
	return res
}

// GroupByKeyTypeValueType returns a mapping of KeyType to []ValueType based on the return
// value of function `f`
//
// Example:
//   groups := GroupByKeyTypeValueType(func(val *ValueType) KeyType {
// 	   return KeyTypeIsh
//   }, vals)
func GroupByKeyTypeValueType(f func(act ValueType) KeyType, coll []ValueType) map[KeyType][]ValueType {
	res := make(map[KeyType][]ValueType)
	EveryValueType(func(act ValueType) bool {
		key := f(act)
		res[key] = append(res[key], act)
		return true
	}, coll)
	return res
}
