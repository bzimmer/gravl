package antonmedv_test

import (
	"context"
	"testing"
	"time"

	"github.com/martinlindhe/unit"
	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/analysis/eval/antonmedv"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

var acts = []*strava.Activity{
	{ID: 1, Type: "Hike", Distance: 100000, ElevationGain: 30, StartDateLocal: time.Date(2009, time.November, 10, 8, 0, 0, 0, time.UTC)},
	{ID: 2, Type: "Ride", Distance: 200000, ElevationGain: 60, StartDateLocal: time.Date(2010, time.December, 10, 8, 0, 0, 0, time.UTC)},
	{ID: 3, Type: "Ride", Distance: 300000, ElevationGain: 90, StartDateLocal: time.Date(2009, time.January, 10, 8, 0, 0, 0, time.UTC)},
	{ID: 4, Type: "Hike", Distance: 400000, ElevationGain: 120, StartDateLocal: time.Date(2010, time.March, 10, 8, 0, 0, 0, time.UTC)},
	{ID: 5, Type: "Ride", Distance: 500000, ElevationGain: 150, StartDateLocal: time.Date(2009, time.April, 10, 8, 0, 0, 0, time.UTC)},
	{ID: 6, Type: "Run", Distance: 600000, ElevationGain: 180, StartDateLocal: time.Date(2011, time.May, 10, 8, 0, 0, 0, time.UTC)},
}

func TestFilterer(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	a.Equal(6, len(acts))

	q, err := antonmedv.Filterer(`.Type == "Ride"`)
	a.NoError(err)
	vals, err := q.Filter(context.Background(), acts)
	a.NotNil(vals)
	a.NoError(err)
	a.Equal(3, len(vals))

	q, err = antonmedv.Filterer(`.Type == "Ride" && .StartDateLocal.Year() == 2010`)
	a.NoError(err)
	vals, err = q.Filter(context.Background(), acts)
	a.NotNil(vals)
	a.NoError(err)
	a.Equal(1, len(vals))
}

func TestMapper(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	a.Equal(6, len(acts))

	q, err := antonmedv.Mapper(`.Type`)
	a.NoError(err)
	vals, err := q.Map(context.Background(), acts)
	a.NotNil(vals)
	a.NoError(err)
	a.Equal(6, len(vals))
	a.Equal("Hike", vals[0])
	a.Equal("Ride", vals[4])
}

func TestEvaluator(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	a.Equal(6, len(acts))

	q, err := antonmedv.Evaluator(`.Type == 'Hike'`)
	a.NoError(err)
	val, err := q.Bool(context.Background(), acts[0])
	a.NoError(err)
	a.True(val)

	q, err = antonmedv.Evaluator(`.Type`)
	a.NoError(err)
	val, err = q.Bool(context.Background(), acts[0])
	a.Error(err)
	a.False(val)

	q, err = antonmedv.Evaluator(`.Type`)
	a.NoError(err)
	yal, err := q.Eval(context.Background(), acts[0])
	a.NoError(err)
	a.Equal("Hike", yal)

	q, err = antonmedv.Evaluator(`[.Type, .Distance]`)
	a.NoError(err)
	yal, err = q.Eval(context.Background(), acts[0])
	a.NoError(err)
	a.Equal([]interface{}{"Hike", unit.Length(100000)}, yal)
}
