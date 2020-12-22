package eddington

import (
	"context"
	"flag"

	"github.com/bzimmer/gravl/pkg/common/stats"
	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/analysis"
)

const Doc = ``

type Eddington struct {
	Units analysis.Units
}

func (a *Eddington) Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
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
	return stats.EddingtonNumber(vals), nil
}

func New() *analysis.Analyzer {
	e := &Eddington{
		Units: analysis.Imperial,
	}
	fs := flag.NewFlagSet("eddington", flag.ExitOnError)
	fs.Var(&analysis.UnitsFlag{Units: &e.Units}, "units", "units to use (default: imperial)")
	return &analysis.Analyzer{
		Name:  "eddington",
		Doc:   Doc,
		Flags: fs,
		Run:   e.Run,
	}
}
