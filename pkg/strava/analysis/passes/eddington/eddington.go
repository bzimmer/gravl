package eddington

import (
	"context"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/analysis"
)

const Doc = `eddington returns the Eddington number for all activities

The Eddington is the largest integer E, where you have cycled at least
E miles (or kilometers) on at least E days.`

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
	return Number(vals), nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "eddington",
		Doc:  Doc,
		Run:  Run,
	}
}
