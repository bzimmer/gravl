package ageride_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/activity/strava"
	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/analysis/passes/ageride"
)

func TestAgeride(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	any := ageride.New()
	a.NotNil(any)
	err := any.Flags.Parse([]string{"--birthday", "2002-07-01"})
	a.NoError(err)

	pass := &analysis.Pass{
		Activities: []*strava.Activity{
			{Type: "Ride", Distance: 34000, ElevationGain: 30, StartDateLocal: time.Date(2020, time.December, 26, 8, 0, 0, 0, time.UTC)},
			{Type: "Ride", Distance: 15500, ElevationGain: 23, StartDateLocal: time.Date(2020, time.December, 27, 8, 0, 0, 0, time.UTC)},
			{Type: "Ride", Distance: 32000, ElevationGain: 85, StartDateLocal: time.Date(2020, time.December, 28, 8, 0, 0, 0, time.UTC)},
			{Type: "Ride", Distance: 17500, ElevationGain: 100, StartDateLocal: time.Date(2020, time.December, 5, 8, 0, 0, 0, time.UTC)},
		},
		Units: analysis.Imperial,
	}
	res, err := any.Run(context.Background(), pass)
	a.NoError(err)
	a.NotNil(res)
	r, ok := res.(*ageride.Result)
	a.True(ok)
	a.NotNil(r)
	a.Equal(2, r.Count)
	a.InDelta(19.8, r.DistanceMedian, 0.1)
	a.InDelta(41.0, r.DistanceTotal, 0.1)
}
