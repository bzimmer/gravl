package stats_test

import (
	"testing"

	"github.com/bzimmer/gravl/pkg/common/stats"
	"github.com/stretchr/testify/assert"
)

func TestEddington(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var r []int
	for i := 0; i < len(rides); i++ {
		r = append(r, int(rides[i]))
	}
	e := stats.EddingtonNumber(r)
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
