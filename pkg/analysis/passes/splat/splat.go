package splat

import (
	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

const doc = `splat returns all activities in the units specified

This analyzer is useful for debugging the filter`

func run(ctx *analysis.Context, pass []*strava.Activity) (interface{}, error) {
	var res []*analysis.Activity
	for i := 0; i < len(pass); i++ {
		act := pass[i]
		res = append(res, analysis.ToActivity(act, ctx.Units))
	}
	return res, nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "splat",
		Doc:  doc,
		Run:  run,
	}
}
