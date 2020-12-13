package stats

import (
	"github.com/logic-building/functional-go/fp"
	"gonum.org/v1/gonum/stat"
)

// https://en.wikipedia.org/wiki/Benford%27s_law
var benfordLaw = []float64{
	0.301, // 1
	0.176, // 2
	0.125, // 3
	0.097, // 4
	0.079, // 5
	0.067, // 6
	0.058, // 7
	0.051, // 8
	0.046, // 9
}

var benfordCount = len(benfordLaw)

type Benford struct {
	Distribution []float64
	ChiSquared   float64
}

func distribution(occ []int) []float64 {
	res := make([]float64, benfordCount)
	sum := float64(fp.ReduceInt(func(x, y int) int {
		return x + y
	}, occ))
	for i := 0; i < len(res); i++ {
		res[i] = float64(occ[i]) / sum
	}
	return res
}

func occurrences(vals []int) []int {
	res := make([]int, benfordCount)
	fp.EveryInt(func(v int) bool {
		if v > 0 {
			res[v-1] = res[v-1] + 1
		}
		return true
	}, fp.MapInt(func(v int) int {
		for v >= 10 {
			v = v / 10
		}
		return v
	}, vals))
	return res
}

func BenfordsLaw(vals []int) Benford {
	dis := distribution(occurrences(vals))
	chi := stat.ChiSquare(dis, benfordLaw)
	return Benford{Distribution: dis, ChiSquared: chi}
}
