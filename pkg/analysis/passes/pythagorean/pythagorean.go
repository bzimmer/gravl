package pythagorean

import (
	"context"
	"math"
	"sort"

	"github.com/bzimmer/gravl/pkg/analysis"
)

type Results struct {
	Activity *analysis.Activity `json:"activity"`
	Number   int                `json:"number"`
}

const doc = ``

func run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	res := make([]*Results, len(pass.Activities))
	for i := 0; i < len(pass.Activities); i++ {
		act := pass.Activities[i]
		dst := act.Distance.Meters()
		elv := act.ElevationGain.Meters()
		n := int(math.Sqrt(math.Pow(dst, 2) + math.Pow(elv, 2)))
		res[i] = &Results{Activity: analysis.ToActivity(act, analysis.Metric), Number: n}
		i++
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Number > res[j].Number
	})
	if len(res) == 0 {
		return nil, nil
	}
	return res[0], nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "pythagorean",
		Doc:  doc,
		Run:  run,
	}
}
