package pythagorean_test

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/activity/strava"
	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/analysis/passes/pythagorean"
)

func TestPythagorean(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	any := pythagorean.New()
	a.NotNil(any)

	pass := &analysis.Pass{
		Activities: []*strava.Activity{
			{ID: 100, Type: "Ride", Distance: 142000, ElevationGain: 240},
			{ID: 200, Type: "Ride", Distance: 152000, ElevationGain: 281},
			{ID: 300, Type: "Ride", Distance: 112000, ElevationGain: 103},
			{ID: 400, Type: "Ride", Distance: 242000, ElevationGain: 220},
			{ID: 500, Type: "Run", Distance: 192000, ElevationGain: 81},
		},
		Units: analysis.Metric,
	}
	res, err := any.Run(context.Background(), pass)
	a.NoError(err)
	a.NotNil(res)
	r, ok := res.(*pythagorean.Results)
	a.True(ok)
	num := int(math.Sqrt(math.Pow(242000, 2) + math.Pow(220, 2)))
	a.Equal(int64(400), r.Activity.ID)
	a.Equal(num, r.Number)
}

func TestPythagoreanEmpty(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	any := pythagorean.New()
	a.NotNil(any)

	pass := &analysis.Pass{Activities: []*strava.Activity{}}
	res, err := any.Run(context.Background(), pass)
	a.NoError(err)
	a.NotNil(res)
	r, ok := res.(*pythagorean.Results)
	a.True(ok)
	a.Equal(0, r.Number)
	a.Nil(r.Activity)
}
