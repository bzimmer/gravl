package cluster_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/analysis/passes/cluster"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

func TestCluster(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	any := cluster.New()
	a.NotNil(any)
	a.NoError(any.Flags.Parse([]string{"--clusters", "2", "--threshold", "0.05"}))

	pass := &analysis.Pass{
		Units: analysis.Imperial,
		Activities: []*strava.Activity{
			{Distance: 100000, ElevationGain: 100},
			{Distance: 100000, ElevationGain: 100},
			{Distance: 100000, ElevationGain: 100},
			{Distance: 400000, ElevationGain: 100},
			{Distance: 100000, ElevationGain: 100},
			{Distance: 400000, ElevationGain: 100},
			{Distance: 100000, ElevationGain: 100},
			{Distance: 400000, ElevationGain: 100},
			{Distance: 100000, ElevationGain: 100},
		},
	}
	res, err := any.Run(context.Background(), pass)
	a.NoError(err)
	a.NotNil(res)
	f, ok := res.([]*cluster.Cluster)
	a.True(ok)
	a.NotNil(f)
	a.Equal(2, len(f))
}
