package climbing

import (
	"context"
	"sort"

	"github.com/bzimmer/gravl/pkg/activity/strava"
	"github.com/bzimmer/gravl/pkg/analysis"
)

const (
	Doc = `All activities exceedinv the Golden Ratio

	https://blog.wahoofitness.com/by-the-numbers-distance-and-elevation/`
	GoldenRatioMetric   = 20
	GoldenRatioImperial = 100
)

type Result struct {
	Activity *analysis.Activity `json:"activity"`
	Number   int                `json:"number"`
}

func Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var climbings []*Result
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		var thr int
		var dst, elv float64
		switch pass.Units {
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
			climbings = append(climbings, &Result{Activity: analysis.ToActivity(act, pass.Units), Number: c})
		}
		return true
	}, pass.Activities)
	sort.Slice(climbings, func(i, j int) bool {
		return climbings[i].Number < climbings[j].Number
	})
	return climbings, nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "climbing",
		Doc:  Doc,
		Run:  Run,
	}
}
