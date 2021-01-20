package rolling

import (
	"flag"
	"fmt"
	"sort"

	"github.com/rs/zerolog/log"
	"gonum.org/v1/gonum/floats"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

const debug = false
const doc = "rolling returns the `window` of activities with the highest accumulated distance."

type Result struct {
	Activities []*analysis.Activity `json:"activities"`
	Distance   float64              `json:"distance"`
}

type rollingWindow struct {
	Window int
}

func (r *rollingWindow) activities(acts []*strava.Activity, idx int, units analysis.Units) []*analysis.Activity {
	var res []*analysis.Activity
	for i := idx; i < idx+r.Window; i++ {
		res = append(res, analysis.ToActivity(acts[i], units))
	}
	return res
}

func (r *rollingWindow) run(ctx *analysis.Context, pass *analysis.Pass) (interface{}, error) {
	if len(pass.Activities) < r.Window {
		log.Warn().Int("n", len(pass.Activities)).Int("window", r.Window).Msg("too few activities")
		return &Result{}, nil
	}
	var dsts = make([]float64, len(pass.Activities))
	var acts = make([]*strava.Activity, len(pass.Activities))
	if n := copy(acts, pass.Activities); n != len(pass.Activities) {
		return nil, fmt.Errorf("%d != %d", n, len(pass.Activities))
	}
	sort.Slice(acts, func(i, j int) bool {
		return acts[i].StartDateLocal.Before(acts[j].StartDateLocal)
	})
	for i := 0; i < len(acts); i++ {
		switch ctx.Units {
		case analysis.Metric:
			dsts[i] = acts[i].Distance.Kilometers()
		case analysis.Imperial:
			dsts[i] = acts[i].Distance.Miles()
		}
	}
	var idx int
	var val float64
	var num = len(dsts) - r.Window
	for i := 0; i <= num; i++ {
		v := floats.Sum(dsts[i : i+r.Window])
		if debug {
			if i == 0 {
				fmt.Println()
			}
			fmt.Printf("%04d,%04d,%04d,%0.3f,%04d,%04d,%0.3f,%0.3f\n",
				len(dsts), num, idx, val, i, i+r.Window, dsts[i], v)
		}
		if v > val {
			idx, val = i, v
			if debug {
				fmt.Printf(" %04d,%03f\n", idx, v)
			}
		}
	}
	res := r.activities(acts, idx, ctx.Units)
	return &Result{Activities: res, Distance: val}, nil
}

func New() *analysis.Analyzer {
	rw := &rollingWindow{Window: 7}
	fs := flag.NewFlagSet("rolling", flag.ExitOnError)
	fs.IntVar(&rw.Window, "window", rw.Window, "the number of days in the window")
	return &analysis.Analyzer{
		Name:  fs.Name(),
		Doc:   doc,
		Flags: fs,
		Run:   rw.run,
	}
}
