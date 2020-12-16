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
		{Distance: 89360, TotalElevationGain: 2447},
		{Distance: 81748, TotalElevationGain: 1756},
		{Distance: 71515, TotalElevationGain: 1740},
		{Distance: 78220, TotalElevationGain: 1720},
		{Distance: 55772, TotalElevationGain: 1633},
	}
	cn := stats.MetricStats.ClimbingNumber(acts)
	a.Equal(5, len(cn))

	acts[1].TotalElevationGain = 1000
	cn = stats.MetricStats.ClimbingNumber(acts)
	a.Equal(4, len(cn))

	cn = stats.ImperialStats.ClimbingNumber(acts)
	a.Equal(4, len(cn))
}
