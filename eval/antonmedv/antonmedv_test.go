package antonmedv_test

import (
	"context"
	"testing"
	"time"

	"github.com/bzimmer/activity/strava"
	"github.com/martinlindhe/unit"
	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/eval/antonmedv"
)

func TestInvalidExpression(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	for _, expr := range []string{".Typ == 'Hike'", ""} {
		f, err := antonmedv.Filterer(expr)
		a.Nil(f)
		a.Error(err)
		m, err := antonmedv.Mapper(expr)
		a.Nil(m)
		a.Error(err)
		e, err := antonmedv.Evaluator(expr)
		a.Nil(e)
		a.Error(err)
	}
}

func activities() []*strava.Activity {
	return []*strava.Activity{
		{ID: 1, Type: "Hike", Distance: 100000, ElevationGain: 30, StartDateLocal: time.Date(2009, time.November, 10, 8, 0, 0, 0, time.UTC)},
		{ID: 2, Type: "Ride", Distance: 200000, ElevationGain: 60, StartDateLocal: time.Date(2010, time.December, 10, 8, 0, 0, 0, time.UTC)},
		{ID: 3, Type: "Ride", Distance: 300000, ElevationGain: 90, StartDateLocal: time.Date(2009, time.January, 10, 8, 0, 0, 0, time.UTC)},
		{ID: 4, Type: "Hike", Distance: 400000, ElevationGain: 120, StartDateLocal: time.Date(2010, time.March, 10, 8, 0, 0, 0, time.UTC)},
		{ID: 5, Type: "Ride", Distance: 500000, ElevationGain: 150, StartDateLocal: time.Date(2009, time.April, 10, 8, 0, 0, 0, time.UTC)},
		{ID: 6, Type: "Run", Distance: 600000, ElevationGain: 180, StartDateLocal: time.Date(2011, time.May, 10, 8, 0, 0, 0, time.UTC)},
	}
}

func TestUserFunctions(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	acts := activities()
	a.Equal(6, len(acts))

	q, err := antonmedv.Mapper("isoweek(.StartDateLocal)")
	a.NoError(err)
	vals, err := q.Map(context.Background(), acts)
	a.NotNil(vals)
	a.NoError(err)
	a.Equal(6, len(vals))
	a.Equal(antonmedv.ISOWeek{Year: 2009, Week: 2}, vals[2])
	a.Equal("[2009 02]", antonmedv.ISOWeek{Year: 2009, Week: 2}.String())

	act := &strava.Activity{ID: 100, Type: "Hike", AverageTemperature: 1.3}
	v, err := antonmedv.Evaluator("F(.AverageTemperature)")
	a.NoError(err)
	u, err := v.Eval(context.Background(), act)
	a.NotNil(u)
	a.NoError(err)
	a.InEpsilon(unit.FromCelsius(1.3).Fahrenheit(), u, 0.1)
}

func TestFilterer(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	acts := activities()
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

	acts := activities()
	a.Equal(6, len(acts))

	q, err := antonmedv.Mapper(".Type")
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

	acts := activities()
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
	a.Equal([]any{"Hike", unit.Length(100000)}, yal)
}
