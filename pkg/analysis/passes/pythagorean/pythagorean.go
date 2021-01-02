package pythagorean

import (
	"context"
	"math"
	"sort"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/strava"
)

type Results struct {
	Activity *analysis.Activity `json:"activity"`
	Number   int                `json:"number"`
}

const Doc = ``

func Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var i int
	res := make([]*Results, len(pass.Activities))
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		dst := act.Distance.Meters()
		elv := act.ElevationGain.Meters()
		n := int(math.Sqrt(math.Pow(dst, 2) + math.Pow(elv, 2)))
		res[i] = &Results{Activity: analysis.ToActivity(act, analysis.Metric), Number: n}
		i++
		return true
	}, pass.Activities)
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
		Doc:  Doc,
		Run:  Run,
	}
}
