package analysis_test

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
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
	q, err := p.Filter("{.Distance.Kilometers() > 150}")
	a.NoError(err)
	a.Equal(2, len(q.Activities))

	q, err = p.Filter("{.Distance.Kilometers() > 150 && .ElevationGain.Meters() > 80}")
	a.NoError(err)
	a.Equal(1, len(q.Activities))

	q, err = p.Filter("{.Distance.Kilometers() < 150 && .ElevationGain.Meters() > 80}")
	a.NoError(err)
	a.Equal(0, len(q.Activities))
}

func TestGroupBy(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	/*
		2009:
		 Hike: [1]
		 Ride: [3, 5]
		2010:
		 Hike: [4]
		 Ride: [2]
		2011:
		 Run: [6]
	*/
	p := &analysis.Pass{
		Activities: []*strava.Activity{
			{ID: 1, Type: "Hike", Distance: 100000, ElevationGain: 30, StartDateLocal: time.Date(2009, time.November, 10, 8, 0, 0, 0, time.UTC)},
			{ID: 2, Type: "Ride", Distance: 200000, ElevationGain: 60, StartDateLocal: time.Date(2010, time.December, 10, 8, 0, 0, 0, time.UTC)},
			{ID: 3, Type: "Ride", Distance: 300000, ElevationGain: 90, StartDateLocal: time.Date(2009, time.January, 10, 8, 0, 0, 0, time.UTC)},
			{ID: 4, Type: "Hike", Distance: 400000, ElevationGain: 120, StartDateLocal: time.Date(2010, time.March, 10, 8, 0, 0, 0, time.UTC)},
			{ID: 5, Type: "Ride", Distance: 500000, ElevationGain: 150, StartDateLocal: time.Date(2009, time.April, 10, 8, 0, 0, 0, time.UTC)},
			{ID: 6, Type: "Run", Distance: 600000, ElevationGain: 180, StartDateLocal: time.Date(2011, time.May, 10, 8, 0, 0, 0, time.UTC)},
		},
	}
	q, err := p.GroupBy()
	a.NoError(err)
	a.NotNil(q)
	a.Equal(0, len(q.Groups))
	a.Equal(6, len(q.Pass.Activities))

	q, err = p.GroupBy("{.Type}")
	a.NoError(err)
	a.NotNil(q)
	a.Equal(3, len(q.Groups))

	q, err = p.GroupBy("{.StartDateLocal.Year()}", "{.Type}")
	a.NoError(err)
	a.NotNil(q)
	a.Equal(3, len(q.Groups))
	sort.Slice(q.Groups, func(i, j int) bool {
		return q.Groups[i].Key < q.Groups[j].Key
	})

	y2009 := q.Groups[0]
	a.Equal(3, len(y2009.Pass.Activities))
	a.Equal(2, len(y2009.Groups))
	y2010 := q.Groups[1]
	a.Equal(2, len(y2010.Pass.Activities))
	a.Equal(2, len(y2010.Groups))
	y2011 := q.Groups[2]
	a.Equal(1, len(y2011.Pass.Activities))
	a.Equal(1, len(y2011.Groups))
}
