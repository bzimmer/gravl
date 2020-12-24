package splat

import (
	"context"
	"flag"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/analysis"
)

const Doc = `splat simply returns all activities in the units specified

This analyzer is useful for debugging the filter.`

type Splat struct {
	Units analysis.Units
}

func (s *Splat) Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var res []*analysis.Activity
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		res = append(res, analysis.ToActivityWithUnits(act, s.Units))
		return true
	}, pass.Activities)
	return res, nil
}

func New() *analysis.Analyzer {
	s := &Splat{
		Units: analysis.Imperial,
	}
	fs := flag.NewFlagSet("splat", flag.ExitOnError)
	fs.Var(&analysis.UnitsFlag{Units: &s.Units}, "units", "units to use")
	return &analysis.Analyzer{
		Name:  fs.Name(),
		Doc:   Doc,
		Flags: fs,
		Run:   s.Run,
	}
}
