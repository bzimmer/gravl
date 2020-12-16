package strava_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/strava"
)

func TestGroupByIntActivityPtr(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	acts := []*strava.Activity{
		{Distance: 10.1}, {Distance: 32.2},
		{Distance: 30.5}, {Distance: 120.9},
	}
	m := strava.GroupByIntActivityPtr(
		func(a *strava.Activity) int {
			return int(math.Mod(a.Distance, 3))
		}, acts)
	a.Equal(3, len(m))
	a.Equal(1, len(m[1]))
	a.Equal(1, len(m[2]))
	a.Equal(2, len(m[0]))
}

func TestFilterActivityPtr(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	acts := []*strava.Activity{
		{Distance: 10.1}, {Distance: 32.2},
		{Distance: 30.5}, {Distance: 120.9},
	}

	m := strava.FilterActivityPtr(nil, acts)
	a.Equal(0, len(m))

	m = strava.FilterActivityPtr(
		func(a *strava.Activity) bool {
			return a.Distance > 31.0
		}, acts)
	a.Equal(2, len(m))
}

func TestMapActivityPtr(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	acts := []*strava.Activity{
		{Distance: 10.0}, {Distance: 32.0},
		{Distance: 30.0}, {Distance: 120.0},
	}
	a.Equal(32.0, acts[1].Distance)

	m := strava.MapActivityPtr(nil, acts)
	a.Equal(0, len(m))

	m = strava.MapActivityPtr(
		func(a *strava.Activity) *strava.Activity {
			a.Distance = a.Distance * 10
			return a
		}, acts)
	a.Equal(4, len(m))
	a.Equal(320.0, acts[1].Distance)
	a.Equal(320.0, m[1].Distance)
}

func TestReduceActivityPtr(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	acts := []*strava.Activity{
		{Distance: 10.0}, {Distance: 32.0},
		{Distance: 30.0}, {Distance: 120.0},
	}
	m := strava.ReduceActivityPtr(
		func(a, b *strava.Activity) *strava.Activity {
			a.Distance = a.Distance + b.Distance
			return a
		}, acts)
	a.Equal(10.0+32.0+30.0+120.0, m.Distance)

	acts = []*strava.Activity{
		{Distance: 10.0}, {Distance: 32.0},
		{Distance: 30.0}, {Distance: 120.0},
	}
	m = strava.ReduceActivityPtr(
		func(a, b *strava.Activity) *strava.Activity {
			a.Distance = a.Distance + b.Distance
			return a
		}, acts, strava.Activity{Distance: 100.0})
	a.Equal(10.0+32.0+30.0+120.0+100.0, m.Distance)

	x := map[string][]*strava.Activity{
		"empty": {},
		"nil":   nil,
	}
	for k, v := range x {
		t.Run(k, func(coll []*strava.Activity) func(*testing.T) {
			return func(t *testing.T) {
				m = strava.ReduceActivityPtr(
					func(a, b *strava.Activity) *strava.Activity {
						a.Distance = a.Distance + b.Distance
						return a
					}, coll)
				a.Nil(m)
			}
		}(v))
	}
}
