package benford

import (
	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

const doc = `benford returns the benford distribution of all the activities`

func run(ctx *analysis.Context, pass []*strava.Activity) (interface{}, error) {
	var vals []int
	for i := 0; i < len(pass); i++ {
		act := pass[i]
		var dst float64
		switch ctx.Units {
		case analysis.Metric:
			dst = act.Distance.Kilometers()
		case analysis.Imperial:
			dst = act.Distance.Miles()
		}
		vals = append(vals, int(dst))
	}
	return Law(vals), nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "benford",
		Doc:  doc,
		Run:  run,
	}
}
