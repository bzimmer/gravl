package climbing

import (
	"context"
	"sort"

	"github.com/bzimmer/gravl/pkg/analysis"
)

const (
	doc = `climbing returns all activities exceeding the Golden Ratio

	https://blog.wahoofitness.com/by-the-numbers-distance-and-elevation/`

	// GoldenRatioMetric in metric units
	GoldenRatioMetric = 20
	// GoldenRatioImperial in imperial units
	GoldenRatioImperial = 100
)

type Result struct {
	Activity *analysis.Activity `json:"activity"`
	Number   int                `json:"number"`
}

func run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var climbings []*Result
	for i := 0; i < len(pass.Activities); i++ {
		act := pass.Activities[i]
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
	}
	sort.Slice(climbings, func(i, j int) bool {
		return climbings[i].Number < climbings[j].Number
	})
	return climbings, nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "climbing",
		Doc:  doc,
		Run:  run,
	}
}
