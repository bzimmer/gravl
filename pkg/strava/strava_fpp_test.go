package strava_test

import (
	"math"
	"testing"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/stretchr/testify/assert"
)

func TestGroupByIntActivityPtr(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	acts := []*strava.Activity{
		{Distance: 10.1},
		{Distance: 32.2},
		{Distance: 30.5},
		{Distance: 120.9},
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
