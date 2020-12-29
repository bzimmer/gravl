package splat

import (
	"context"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/analysis"
)

const Doc = `splat simply returns all activities in the units specified

This analyzer is useful for debugging the filter.`

func Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var res []*analysis.Activity
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		res = append(res, analysis.ToActivity(act, pass.Units))
		return true
	}, pass.Activities)
	return res, nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "splat",
		Doc:  Doc,
		Run:  Run,
	}
}
