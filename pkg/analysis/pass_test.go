package analysis_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/strava"
)

func TestFilter(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	p := &analysis.Pass{
		Activities: []*strava.Activity{
			{Distance: 142000, ElevationGain: 30},
			{Distance: 155000, ElevationGain: 23},
			{Distance: 202000, ElevationGain: 85},
		},
	}
	q, err := p.FilterExpr("{.Distance.Kilometers() > 150}")
	a.NoError(err)
	a.Equal(2, len(q.Activities))

	q, err = p.FilterExpr("{.Distance.Kilometers() > 150 && .ElevationGain.Meters() > 80}")
	a.NoError(err)
	a.Equal(1, len(q.Activities))

	q, err = p.FilterExpr("{.Distance.Kilometers() < 150 && .ElevationGain.Meters() > 80}")
	a.NoError(err)
	a.Equal(0, len(q.Activities))
}

func TestGroupBy(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	p := &analysis.Pass{
		Activities: []*strava.Activity{
			{Type: "Hike", Distance: 142000, ElevationGain: 30},
			{Type: "Ride", Distance: 155000, ElevationGain: 23},
			{Type: "Ride", Distance: 302000, ElevationGain: 120},
			{Type: "Hike", Distance: 240200, ElevationGain: 232},
		},
	}
	q, err := p.GroupByExpr("{.Type}")
	a.NoError(err)
	a.Equal(2, len(q))
	a.Equal(2, len(q["Hike"].Activities))
	a.Equal(2, len(q["Ride"].Activities))
}
