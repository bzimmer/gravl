package hourrecord

// https://swinny.net/Strava/-4745-My-Hour-Record

import (
	"context"
	"sort"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

const doc = `hourrecord returns the longest distance traveled (in miles | kilometers) exceeding the average speed (mph | mps).`

func run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var res []*strava.Activity
	for i := 0; i < len(pass.Activities); i++ {
		act := pass.Activities[i]
		dst := act.Distance
		spd := act.AverageSpeed
		if float64(dst) >= float64(spd) {
			res = append(res, act)
		}
	}
	if len(res) == 0 {
		return nil, nil
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].AverageSpeed > res[j].AverageSpeed
	})
	return analysis.ToActivity(res[0], pass.Units), nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "hourrecord",
		Doc:  doc,
		Run:  run,
	}
}
