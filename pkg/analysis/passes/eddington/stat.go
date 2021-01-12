package eddington

type Eddington struct {
	Number     int         `json:"number"`
	Numbers    []int       `json:"numbers"`
	Motivation map[int]int `json:"motivation"`
}

// Number computes the Eddington number from a series of rides
func Number(rides []int) *Eddington {
	n, above := len(rides), 0
	e := &Eddington{
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
