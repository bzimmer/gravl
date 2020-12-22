package festive500

import (
	"context"
	"time"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/analysis"
)

const Doc = ``

type Result struct {
	Activities []*analysis.Activity `json:"activities"`
	Distance   float64              `json:"distance"`
}

func Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var dst float64
	var res []*analysis.Activity
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		_, month, date := act.StartDateLocal.Date()
		ok := (month == time.December && date >= 24 && date <= 31)
		if ok {
			dst = dst + act.Distance.Kilometers()
			res = append(res, analysis.ToActivityWithUnits(act, analysis.Metric))
		}
		return true
	}, pass.Activities)
	return &Result{Activities: res, Distance: dst}, nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "festive500",
		Doc:  Doc,
		Run:  Run,
	}
}
