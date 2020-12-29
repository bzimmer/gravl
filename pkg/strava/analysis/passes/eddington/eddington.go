package eddington

import (
	"context"
	"flag"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/analysis"
)

const Doc = `eddington returns the Eddington number for all activities

The Eddington is the largest integer E, where you have cycled at least
E miles (or kilometers) on at least E days.`

type E struct {
	Units analysis.Units
}

func (a *E) Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var vals []int
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		var dst float64
		switch a.Units {
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
	e := &E{
		Units: analysis.Imperial,
	}
	fs := flag.NewFlagSet("eddington", flag.ExitOnError)
	fs.Var(&analysis.UnitsFlag{Units: &e.Units}, "units", "units to use")
	return &analysis.Analyzer{
		Name:  fs.Name(),
		Doc:   Doc,
		Flags: fs,
		Run:   e.Run,
	}
}
