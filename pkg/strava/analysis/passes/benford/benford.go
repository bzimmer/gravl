package benford

import (
	"context"
	"flag"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/analysis"
)

const Doc = ``

type B struct {
	Units analysis.Units
}

func (a *B) Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
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
	return Law(vals), nil
}

func New() *analysis.Analyzer {
	e := &B{
		Units: analysis.Imperial,
	}
	fs := flag.NewFlagSet("benford", flag.ExitOnError)
	fs.Var(&analysis.UnitsFlag{Units: &e.Units}, "units", "units to use")
	return &analysis.Analyzer{
		Name:  fs.Name(),
		Doc:   Doc,
		Flags: fs,
		Run:   e.Run,
	}
}
