package festive500_test

import (
	"context"
	"testing"
	"time"

	"github.com/martinlindhe/unit"
	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/analysis/passes/festive500"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

func TestFestive500(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	any := festive500.New()
	a.NotNil(any)

	pass := &analysis.Pass{
		Activities: []*strava.Activity{
			{Type: "Ride", Distance: 142000, ElevationGain: 30, StartDateLocal: time.Date(2020, time.December, 26, 8, 0, 0, 0, time.UTC)},
			{Type: "Ride", Distance: 155000, ElevationGain: 23, StartDateLocal: time.Date(2020, time.December, 27, 8, 0, 0, 0, time.UTC)},
			{Type: "Ride", Distance: 202000, ElevationGain: 85, StartDateLocal: time.Date(2020, time.December, 28, 8, 0, 0, 0, time.UTC)},
			{Type: "Ride", Distance: 175000, ElevationGain: 100, StartDateLocal: time.Date(2020, time.December, 5, 8, 0, 0, 0, time.UTC)},
		},
	}
	res, err := any.Run(context.Background(), pass)
	a.NoError(err)
	a.NotNil(res)
	f, ok := res.(*festive500.Result)
	a.True(ok)
	a.Equal(unit.Length(142000+155000+202000).Kilometers(), f.DistanceCompleted)
}
