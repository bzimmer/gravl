package benford

import (
	"context"

	"github.com/bzimmer/gravl/pkg/analysis"
)

const doc = ``

func run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var vals []int
	for i := 0; i < len(pass.Activities); i++ {
		act := pass.Activities[i]
		var dst float64
		switch pass.Units {
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
