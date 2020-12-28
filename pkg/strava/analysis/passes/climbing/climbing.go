package climbing

import (
	"context"
	"flag"
	"sort"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/analysis"
)

const (
	Doc = `All activities exceedinv the Golden Ratio

	https://blog.wahoofitness.com/by-the-numbers-distance-and-elevation/`
	GoldenRatioMetric   = 30
	GoldenRatioImperial = 100
)

type Result struct {
	Activity *analysis.Activity `json:"activity"`
	Number   int                `json:"number"`
}

type Climbing struct {
	Units analysis.Units
}

func (a *Climbing) Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var climbings []*Result
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		var thr int
		var dst, elv float64
		switch a.Units {
		case analysis.Metric:
			thr = GoldenRatioMetric
			dst = act.Distance.Kilometers()
			elv = act.ElevationGain.Meters()
		case analysis.Imperial:
			thr = GoldenRatioImperial
			dst = act.Distance.Miles()
			elv = act.ElevationGain.Feet()
		}
		c := int(elv / dst)
		if c > thr {
			climbings = append(climbings, &Result{Activity: analysis.ToActivity(act, a.Units), Number: c})
		}
		return true
	}, pass.Activities)
	sort.Slice(climbings, func(i, j int) bool {
		return climbings[i].Number < climbings[j].Number
	})
	return climbings, nil
}

func New() *analysis.Analyzer {
	c := &Climbing{
		Units: analysis.Imperial,
	}
	fs := flag.NewFlagSet("climbing", flag.ExitOnError)
	fs.Var(&analysis.UnitsFlag{Units: &c.Units}, "units", "units to use")
	return &analysis.Analyzer{
		Name:  fs.Name(),
		Doc:   Doc,
		Flags: fs,
		Run:   c.Run,
	}
}
