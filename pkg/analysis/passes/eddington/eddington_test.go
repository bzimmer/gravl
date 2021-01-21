package eddington_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/analysis/passes/eddington"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

func TestEddingtonAnalysisEmpty(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	any := eddington.New()
	a.NotNil(any)

	ctx := analysis.WithContext(context.Background(), analysis.Metric)
	res, err := any.Run(ctx, []*strava.Activity{})
	a.NoError(err)
	a.NotNil(res)
	r, ok := res.(*eddington.Eddington)
	a.True(ok)
	a.Equal(0, r.Number)
	a.Equal([]int{}, r.Numbers)
	a.Equal(map[int]int{}, r.Motivation)
}

func TestEddingtonAnalysis(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	any := eddington.New()
	a.NotNil(any)

	pass := []*strava.Activity{
		{ID: 100, Type: "Ride", Distance: 1000, ElevationGain: 240},
		{ID: 200, Type: "Ride", Distance: 2000, ElevationGain: 281},
		{ID: 300, Type: "Ride", Distance: 1000, ElevationGain: 103},
		{ID: 400, Type: "Ride", Distance: 3000, ElevationGain: 220},
		{ID: 500, Type: "Ride", Distance: 2000, ElevationGain: 101},
		{ID: 600, Type: "Ride", Distance: 1000, ElevationGain: 220},
	}
	ctx := analysis.WithContext(context.Background(), analysis.Metric)
	res, err := any.Run(ctx, pass)
	a.NoError(err)
	a.NotNil(res)
	r, ok := res.(*eddington.Eddington)
	a.True(ok)
	a.Equal(2, r.Number)
	a.Equal([]int{1, 1, 1, 2, 2, 2}, r.Numbers)
	a.Equal(map[int]int{3: 1}, r.Motivation)
}

func TestEddingtonNumber(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var r []int
	for i := 0; i < len(rides); i++ {
		r = append(r, int(rides[i]))
	}
	e := eddington.Number(r)
	a.Equal(21, e.Number)
}

var rides = []float64{
	5.43,
	5.414,
	32.198,
	30.322,
	18.117,
	145.352,
	22.967,
	29.585,
	29.939,
	157.036,
	24.946,
	25.303,
	51.146,
	23.944,
	6.01,
	24.4,
	30.903,
	39.48,
	5.907,
	35.825,
	6.768,
	71.515,
	7.494,
	32.614,
	23.183,
	17.455,
	135.918,
	6.577,
	27.225,
	22.061,
}
