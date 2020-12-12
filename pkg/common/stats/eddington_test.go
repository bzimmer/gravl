package stats_test

import (
	"testing"

	"github.com/bzimmer/gravl/pkg/common/stats"
	"github.com/stretchr/testify/assert"
)

// (geo.EddingtonNumber) {
//  Number: (int) 21,
//  Numbers: ([]int) (len=30 cap=30) {
//   (int) 1,
//   (int) 2,
//   (int) 3,
//   (int) 4,
//   (int) 5,
//   (int) 5,
//   (int) 5,
//   (int) 6,
//   (int) 7,
//   (int) 8,
//   (int) 9,
//   (int) 10,
//   (int) 11,
//   (int) 12,
//   (int) 12,
//   (int) 13,
//   (int) 14,
//   (int) 15,
//   (int) 15,
//   (int) 16,
//   (int) 16,
//   (int) 17,
//   (int) 17,
//   (int) 18,
//   (int) 18,
//   (int) 18,
//   (int) 19,
//   (int) 19,
//   (int) 20,
//   (int) 21
//  },
//  Motivation: (map[int]int) (len=6) {
//   (int) 29: (int) 2,
//   (int) 24: (int) 2,
//   (int) 25: (int) 1,
//   (int) 23: (int) 2,
//   (int) 22: (int) 2,
//   (int) 27: (int) 1
//  },
// }

func TestEddington(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var r []int
	for i := 0; i < len(rides); i++ {
		r = append(r, int(rides[i]))
	}
	e := stats.Eddington(r)
	a.Equal(21, e.Number)
}

var rides = []float64{
	5.43,
	5.414,
	32.198,
	30.322,
	18.117,
	145.352,
	22.967,
	29.585,
	29.939,
	157.036,
	24.946,
	25.303,
	51.146,
	23.944,
	6.01,
	24.4,
	30.903,
	39.48,
	5.907,
	35.825,
	6.768,
	71.515,
	7.494,
	32.614,
	23.183,
	17.455,
	135.918,
	6.577,
	27.225,
	22.061,
}
