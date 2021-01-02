package hourrecord

import (
	"context"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/strava"
)

const Doc = ``

func Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
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
	return analysis.ToActivity(act, analysis.Imperial), nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "hourrecord",
		Doc:  Doc,
		Run:  Run,
	}
}
