package eddington

import (
	"context"

	"github.com/bzimmer/gravl/pkg/analysis"
)

const doc = `eddington returns the Eddington number for all activities

The Eddington is the largest integer E, where you have cycled at least
E miles (or kilometers) on at least E days.`

func run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var dst float64
	var vals = make([]int, len(pass.Activities))
	for i := 0; i < len(pass.Activities); i++ {
		act := pass.Activities[i]
		switch pass.Units {
		case analysis.Metric:
			dst = act.Distance.Kilometers()
		case analysis.Imperial:
			dst = act.Distance.Miles()
		}
		vals[i] = int(dst)
	}
	return Number(vals), nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "eddington",
		Doc:  doc,
		Run:  run,
	}
}
