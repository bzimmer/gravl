package stats

type Eddington struct {
	Number     int
	Numbers    []int
	Motivation map[int]int
}

// Eddington computes the Eddington number from a series of rides
func EddingtonNumber(rides []int) Eddington {
	n, above := len(rides), 0
	e := Eddington{
		Number:     0,
		Numbers:    make([]int, n),
		Motivation: make(map[int]int),
	}
	for i, ride := range rides {
		if ride > e.Number {
			above++
			if ride < n {
				e.Motivation[ride]++
			}
			if above > e.Number {
				e.Number++
				val, ok := e.Motivation[e.Number]
				if ok {
					above -= val
					delete(e.Motivation, e.Number)
				}
			}
		}
		e.Numbers[i] = e.Number
	}
	return e
}
