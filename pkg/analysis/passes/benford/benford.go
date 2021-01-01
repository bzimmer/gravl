package benford

import (
	"context"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/strava"
)

const Doc = ``

func Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var vals []int
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		var dst float64
		switch pass.Units {
		case analysis.Metric:
			dst = act.Distance.Kilometers()
		case analysis.Imperial:
			dst = act.Distance.Miles()
		}
		vals = append(vals, int(dst))
		return true
	}, pass.Activities)
	return Law(vals), nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "benford",
		Doc:  Doc,
		Run:  Run,
	}
}
