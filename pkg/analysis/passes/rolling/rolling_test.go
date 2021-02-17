package rolling_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/analysis/passes/rolling"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
	"github.com/stretchr/testify/assert"
)

func TestRolling(t *testing.T) {
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
		name   string
		window int
		zero   bool
		value  float64
	}{
		{name: "seven", window: 7, value: 789.95},
		{name: "three", window: 3, value: 362.45},
		{name: "many", window: len(pass) + 1, zero: true},
	}
	for _, tt := range tests {
		v := tt
		t.Run(v.name, func(t *testing.T) {
			any := rolling.New()
			a.NoError(any.Flags.Lookup("window").Value.Set(fmt.Sprintf("%d", v.window)))
			ctx := analysis.WithContext(context.Background(), analysis.Imperial)
			res, err := any.Run(ctx, pass)
			a.NoError(err)
			a.NotNil(res)
			r, ok := res.(*rolling.Result)
			a.True(ok)
			if v.zero {
				a.Equal(0, len(r.Activities))
				a.Equal(0.0, r.Distance)
				return
			}
			a.Equal(v.window, len(r.Activities))
			a.InEpsilon(v.value, r.Distance, 0.1)
		})
	}
}
