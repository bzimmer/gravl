package strava

// GroupByIntActivityPtr returns a mapping int to []*Activity based on the return value of function `f`
//
// Example:
//   groups := strava.GroupByIntActivityPtr(func(act *strava.Activity) int {
// 	  return act.StartDateLocal.Year()
//   }, acts)
func GroupByIntActivityPtr(f func(act *Activity) int, acts []*Activity) map[int][]*Activity {
	res := make(map[int][]*Activity)
	EveryActivityPtr(func(act *Activity) bool {
		key := f(act)
		res[key] = append(res[key], act)
		return true
	}, acts)
	return res
}
