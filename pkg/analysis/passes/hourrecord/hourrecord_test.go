package hourrecord_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/activity/strava"
	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/analysis/passes/hourrecord"
)

func TestHourRecord(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	any := hourrecord.New()
	a.NotNil(any)

	pass := &analysis.Pass{
		Activities: []*strava.Activity{
			{ID: 100, Type: "Ride", Distance: 142000, AverageSpeed: 24},
			{ID: 200, Type: "Ride", Distance: 152000, AverageSpeed: 28},
			{ID: 300, Type: "Ride", Distance: 112000, AverageSpeed: 10},
			{ID: 400, Type: "Ride", Distance: 242000, AverageSpeed: 22},
		},
	}
	res, err := any.Run(context.Background(), pass)
	a.NoError(err)
	a.NotNil(res)
	f, ok := res.(*analysis.Activity)
	a.True(ok)
	a.Equal(int64(200), f.ID)
}
