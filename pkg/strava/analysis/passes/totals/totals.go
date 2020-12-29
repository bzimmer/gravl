package totals

import (
	"context"
	"time"

	"github.com/martinlindhe/unit"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/analysis"
)

const Doc = ``

type Result struct {
	Distance      unit.Length   `json:"distance"`
	ElevationGain unit.Length   `json:"elevation_gain"`
	MovingTime    time.Duration `json:"moving_time"`
}

func Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var dst, elv float64
	var dur time.Duration
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		dur = dur + (time.Duration(act.MovingTime) * time.Second)
		switch pass.Units {
		case analysis.Metric:
			dst = dst + act.Distance.Kilometers()
			elv = elv + act.ElevationGain.Meters()
		case analysis.Imperial:
			dst = dst + act.Distance.Miles()
			elv = elv + act.ElevationGain.Feet()
		}
		return true
	}, pass.Activities)
	return &Result{
		Distance:      unit.Length(dst),
		ElevationGain: unit.Length(elv),
		MovingTime:    dur / (time.Second * 60 * 60),
	}, nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "totals",
		Doc:  Doc,
		Run:  Run,
	}
}
