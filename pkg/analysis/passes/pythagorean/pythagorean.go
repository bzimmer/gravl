package pythagorean

import (
	"math"
	"sort"

	"github.com/bzimmer/gravl/pkg/analysis"
)

type Results struct {
	Activity *analysis.Activity `json:"activity"`
	Number   int                `json:"number"`
}

const doc = `pythagorean results the activity with the highest pythagorean value.

The pythagorean value is the sqrt of the distance^2 + elevation^2.`

func run(ctx *analysis.Context, pass *analysis.Pass) (interface{}, error) {
	if len(pass.Activities) == 0 {
		return &Results{}, nil
	}
	res := make([]*Results, len(pass.Activities))
	for i := 0; i < len(pass.Activities); i++ {
		act := pass.Activities[i]
		dst := act.Distance.Meters()
		elv := act.ElevationGain.Meters()
		n := int(math.Sqrt(math.Pow(dst, 2) + math.Pow(elv, 2)))
		res[i] = &Results{Activity: analysis.ToActivity(act, analysis.Metric), Number: n}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Number > res[j].Number
	})
	return res[0], nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "pythagorean",
		Doc:  doc,
		Run:  run,
	}
}
