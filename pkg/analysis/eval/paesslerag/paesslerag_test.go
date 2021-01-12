package paesslerag_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/analysis/eval/paesslerag"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

func TestFilter(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	v := paesslerag.New()
	acts := []*strava.Activity{
		{ID: 1, Type: "Hike", Distance: 100000, ElevationGain: 30, StartDateLocal: time.Date(2009, time.November, 10, 8, 0, 0, 0, time.UTC)},
		{ID: 2, Type: "Ride", Distance: 200000, ElevationGain: 60, StartDateLocal: time.Date(2010, time.December, 10, 8, 0, 0, 0, time.UTC)},
		{ID: 3, Type: "Ride", Distance: 300000, ElevationGain: 90, StartDateLocal: time.Date(2009, time.January, 10, 8, 0, 0, 0, time.UTC)},
		{ID: 4, Type: "Hike", Distance: 400000, ElevationGain: 120, StartDateLocal: time.Date(2010, time.March, 10, 8, 0, 0, 0, time.UTC)},
		{ID: 5, Type: "Ride", Distance: 500000, ElevationGain: 150, StartDateLocal: time.Date(2009, time.April, 10, 8, 0, 0, 0, time.UTC)},
		{ID: 6, Type: "Run", Distance: 600000, ElevationGain: 180, StartDateLocal: time.Date(2011, time.May, 10, 8, 0, 0, 0, time.UTC)},
	}
	a.Equal(6, len(acts))

	q := `Type == "Ride"`
	acts, err := v.Filter(context.Background(), q, acts)
	a.NotNil(acts)
	a.NoError(err)
	a.Equal(3, len(acts))

	q = `Type == "Ride" && StartDateLocal.Year() == 2010`
	acts, err = v.Filter(context.Background(), q, acts)
	a.NotNil(acts)
	a.NoError(err)
	a.Equal(1, len(acts))
}

func TestGroup(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	v := paesslerag.New()
	q := `Type`
	acts := []*strava.Activity{{Type: "Run"}, {Type: "Ride"}, {Type: "Run"}, {Type: "Walk"}}
	a.Equal(4, len(acts))
	groups, err := v.Group(context.Background(), q, acts)
	a.NotNil(acts)
	a.NoError(err)
	a.Equal(3, len(groups))
}
