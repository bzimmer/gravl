package climbing_test

import (
	"context"
	"testing"

	"github.com/martinlindhe/unit"
	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/analysis/passes/climbing"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

func TestClimbing(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	any := climbing.New()
	a.NotNil(any)

	mi := unit.Length(155000).Miles()
	pass := &analysis.Pass{
		Units: analysis.Imperial,
		Activities: []*strava.Activity{
			{Distance: 142000, ElevationGain: 30},
			{Distance: 155000, ElevationGain: unit.Length((mi * climbing.GoldenRatioImperial) + 10)},
			{Distance: 202000, ElevationGain: 85},
		},
	}
	res, err := any.Run(context.Background(), pass)
	a.NoError(err)
	a.NotNil(res)
	f, ok := res.([]*climbing.Result)
	a.True(ok)
	a.NotNil(f)
	a.Equal(328, f[0].Number)
}
