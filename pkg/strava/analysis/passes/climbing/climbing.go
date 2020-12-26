package climbing

import (
	"context"
	"flag"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/analysis"
)

const Doc = ``

type Result struct {
	Activity *analysis.Activity `json:"activity"`
	Number   int                `json:"number"`
}

type Climbing struct {
	Units     analysis.Units
	Threshold int
}

func (a *Climbing) Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var climbings []*Result
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		var dst, elv float64
		switch a.Units {
		case analysis.Metric:
			dst = act.Distance.Kilometers()
			elv = act.ElevationGain.Meters()
		case analysis.Imperial:
			dst = act.Distance.Miles()
			elv = act.ElevationGain.Feet()
		}
		c := int(elv / dst)
		if c > a.Threshold {
			climbings = append(climbings, &Result{Activity: analysis.ToActivity(act, a.Units), Number: c})
		}
		return true
	}, pass.Activities)
	return climbings, nil
}

func New() *analysis.Analyzer {
	c := &Climbing{
		Units:     analysis.Imperial,
		Threshold: 100,
	}
	fs := flag.NewFlagSet("climbing", flag.ExitOnError)
	fs.IntVar(&c.Threshold, "threshold", c.Threshold, "climbing threshold")
	fs.Var(&analysis.UnitsFlag{Units: &c.Units}, "units", "units to use")
	return &analysis.Analyzer{
		Name:  fs.Name(),
		Doc:   Doc,
		Flags: fs,
		Run:   c.Run,
	}
}
