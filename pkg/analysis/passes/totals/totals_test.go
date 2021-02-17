package totals_test

import (
	"context"
	"testing"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/analysis/passes/totals"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
	"github.com/stretchr/testify/assert"
)

func TestTotals(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	pass := []*strava.Activity{
		{ID: 100, Type: "Ride", Distance: 142000, ElevationGain: 240},
		{ID: 200, Type: "Ride", Distance: 152000, ElevationGain: 281},
		{ID: 300, Type: "Ride", Distance: 112000, ElevationGain: 103},
		{ID: 400, Type: "Ride", Distance: 242000, ElevationGain: 220},
		{ID: 500, Type: "Ride", Distance: 192000, ElevationGain: 81},
		{ID: 600, Type: "Ride", Distance: 142000, ElevationGain: 651},
		{ID: 700, Type: "Ride", Distance: 194000, ElevationGain: 109},
		{ID: 800, Type: "Ride", Distance: 191200, ElevationGain: 223},
		{ID: 900, Type: "Ride", Distance: 198100, ElevationGain: 281},
	}
	tests := []struct {
		name     string
		unit     analysis.Units
		count    int
		distance float64
	}{
		{name: "all-metric", unit: analysis.Metric, count: len(pass), distance: 1565.3},
		{name: "all-imperial", unit: analysis.Imperial, count: len(pass), distance: 972.6},
	}
	for _, tt := range tests {
		v := tt
		t.Run(v.name, func(t *testing.T) {
			ctx := analysis.WithContext(context.Background(), v.unit)
			any := totals.New()
			res, err := any.Run(ctx, pass)
			a.NoError(err)
			a.NotNil(res)
			r, ok := res.(*totals.Result)
			a.True(ok)
			a.Equal(v.count, r.Count)
			a.InEpsilon(v.distance, r.Distance, 0.1)
		})
	}
}
