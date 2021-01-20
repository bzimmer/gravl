package analysis_test

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

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
var acts = []*strava.Activity{
	{ID: 1, Type: "Hike", Distance: 100000, ElevationGain: 30, StartDateLocal: time.Date(2009, time.November, 10, 8, 0, 0, 0, time.UTC)},
	{ID: 2, Type: "Ride", Distance: 200000, ElevationGain: 60, StartDateLocal: time.Date(2010, time.December, 10, 8, 0, 0, 0, time.UTC)},
	{ID: 3, Type: "Ride", Distance: 300000, ElevationGain: 90, StartDateLocal: time.Date(2009, time.January, 10, 8, 0, 0, 0, time.UTC)},
	{ID: 4, Type: "Hike", Distance: 400000, ElevationGain: 120, StartDateLocal: time.Date(2010, time.March, 10, 8, 0, 0, 0, time.UTC)},
	{ID: 5, Type: "Ride", Distance: 500000, ElevationGain: 150, StartDateLocal: time.Date(2009, time.April, 10, 8, 0, 0, 0, time.UTC)},
	{ID: 6, Type: "Run", Distance: 600000, ElevationGain: 180, StartDateLocal: time.Date(2011, time.May, 10, 8, 0, 0, 0, time.UTC)},
}

type mapper struct {
	fn func(*strava.Activity) interface{}
}

// Map over the collection of activities producing a slice of expression evaluation values
func (x *mapper) Map(ctx context.Context, acts []*strava.Activity) ([]interface{}, error) {
	res := make([]interface{}, len(acts))
	for i := range acts {
		res[i] = x.fn(acts[i])
	}
	return res, nil
}

func TestGroup(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	ctx := analysis.WithContext(context.Background(), analysis.Imperial)
	q, err := analysis.Group(ctx, acts)
	a.NoError(err)
	a.NotNil(q)
	a.Equal(0, len(q.Children))
	a.Equal(6, len(q.Activities))

	fnyear := func(act *strava.Activity) interface{} {
		return act.StartDateLocal.Year()
	}
	q, err = analysis.Group(ctx, acts, &mapper{fnyear})
	a.NoError(err)
	a.NotNil(q)
	a.Equal(3, len(q.Children))

	fntype := func(act *strava.Activity) interface{} {
		return act.Type
	}
	q, err = analysis.Group(ctx, acts, &mapper{fnyear}, &mapper{fntype})
	a.NoError(err)
	a.NotNil(q)
	a.Equal(3, len(q.Children))
	sort.Slice(q.Children, func(i, j int) bool {
		return q.Children[i].Key < q.Children[j].Key
	})

	y2009 := q.Children[0]
	a.Equal(3, len(y2009.Activities))
	a.Equal(2, len(y2009.Children))
	y2010 := q.Children[1]
	a.Equal(2, len(y2010.Activities))
	a.Equal(2, len(y2010.Children))
	y2011 := q.Children[2]
	a.Equal(1, len(y2011.Activities))
	a.Equal(1, len(y2011.Children))
}
