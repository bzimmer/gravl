package pythagorean

import (
	"math"
	"sort"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

// Result of running the pythagorean algorithm
type Result struct {
	Activity *analysis.Activity `json:"activity"`
	Number   int                `json:"number"`
}

const doc = `pythagorean determines the activity with the highest pythagorean value
defined as the sqrt(distance^2 + elevation^2) in meters`

func run(ctx *analysis.Context, pass []*strava.Activity) (interface{}, error) {
	if len(pass) == 0 {
		return &Result{}, nil
	}
	res := make([]*Result, len(pass))
	for i := 0; i < len(pass); i++ {
		act := pass[i]
		dst := act.Distance.Meters()
		elv := act.ElevationGain.Meters()
		n := int(math.Sqrt(math.Pow(dst, 2) + math.Pow(elv, 2)))
		res[i] = &Result{Activity: analysis.ToActivity(act, analysis.Metric), Number: n}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Number > res[j].Number
	})
	return res[0], nil
}

// New returns a new analyzer
func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "pythagorean",
		Doc:  doc,
		Run:  run,
	}
}
