package koms_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/activity/strava"
	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/analysis/passes/koms"
)

func TestKOMs(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	any := koms.New()
	a.NotNil(any)

	pass := &analysis.Pass{
		Units: analysis.Imperial,
		Activities: []*strava.Activity{
			{Distance: 142000, ElevationGain: 30,
				SegmentEfforts: []*strava.SegmentEffort{
					{
						Achievements: []*strava.Achievement{
							{Rank: 1, Type: "overall"},
							{Rank: 1, Type: "male"},
						}},
					{
						Achievements: []*strava.Achievement{
							{Rank: 3, Type: "overall"},
							{Rank: 4, Type: "male"},
						}},
				},
			},
			{Distance: 155000, ElevationGain: 50,
				SegmentEfforts: []*strava.SegmentEffort{{
					Achievements: []*strava.Achievement{
						{Rank: 1, Type: "overall"},
						{Rank: 2, Type: "overall"},
					},
				}}},
			{Distance: 202000, ElevationGain: 85},
		},
	}
	res, err := any.Run(context.Background(), pass)
	a.NoError(err)
	a.NotNil(res)
	f, ok := res.([]*strava.SegmentEffort)
	a.True(ok)
	a.NotNil(f)
	a.Equal(2, len(f))
}
