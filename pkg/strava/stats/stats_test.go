package stats_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/stats"
)

func TestClimbingNumber(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	acts := []*strava.Activity{
		{Distance: 89360, ElevationGain: 2447},
		{Distance: 81748, ElevationGain: 1756},
		{Distance: 71515, ElevationGain: 1740},
		{Distance: 78220, ElevationGain: 1720},
		{Distance: 55772, ElevationGain: 1633},
	}
	cn := stats.ClimbingNumber(acts, stats.Metric, 20)
	a.Equal(5, len(cn))

	acts[1].ElevationGain = 1000
	cn = stats.ClimbingNumber(acts, stats.Metric, 20)
	a.Equal(4, len(cn))
}
