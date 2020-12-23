package stats

import (
	"gonum.org/v1/gonum/stat"
)

// https://en.wikipedia.org/wiki/Benford%27s_law
var benfordDistribution = []float64{
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

var benfordCount = len(benfordDistribution)

type Benford struct {
	Distribution []float64 `json:"distribution"`
	ChiSquared   float64   `json:"chi_squared"`
}

func distribution(occ []int) []float64 {
	res := make([]float64, benfordCount)
	sum := 0.0
	for i := 0; i < len(occ); i++ {
		sum = sum + float64(occ[i])
	}
	for i := 0; i < len(res); i++ {
		res[i] = float64(occ[i]) / sum
	}
	return res
}

func occurrences(vals []int) []int {
	res := make([]int, benfordCount)
	for i := 0; i < len(vals); i++ {
		v := vals[i]
		for v >= 10 {
			v = v / 10
		}
		if v > 0 {
			res[v-1] = res[v-1] + 1
		}
	}
	return res
}

func BenfordsLaw(vals []int) Benford {
	dis := distribution(occurrences(vals))
	chi := stat.ChiSquare(dis, benfordDistribution)
	return Benford{Distribution: dis, ChiSquared: chi}
}
