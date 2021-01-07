package hourrecord

// https://swinny.net/Strava/-4745-My-Hour-Record

import (
	"context"

	"github.com/bzimmer/gravl/pkg/activity/strava"
	"github.com/bzimmer/gravl/pkg/analysis"
)

const doc = `The longest distance traveled (in miles | kilometers) exceeding the average speed (mph | mps).`

func run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	act := strava.ReduceActivityPtr(func(act0, act1 *strava.Activity) *strava.Activity {
		if act0.AverageSpeed > act1.AverageSpeed {
			return act0
		}
		return act1
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		dst := act.Distance
		spd := act.AverageSpeed
		return float64(dst) >= float64(spd)
	}, pass.Activities))
	if act == nil {
		return nil, nil
	}
	return analysis.ToActivity(act, pass.Units), nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "hourrecord",
		Doc:  doc,
		Run:  run,
	}
}
