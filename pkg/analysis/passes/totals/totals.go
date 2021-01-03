package totals

import (
	"context"
	"time"

	"github.com/martinlindhe/unit"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/strava"
)

const Doc = ``

type Centuries struct {
	Metric   int `json:"metric"`
	Imperial int `json:"imperial"`
}

type Result struct {
	Count         int           `json:"count"`
	Distance      unit.Length   `json:"distance"`
	ElevationGain unit.Length   `json:"elevation_gain"`
	MovingTime    time.Duration `json:"moving_time"`
	Centuries     Centuries     `json:"centuries"`
}

func Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var cen Centuries
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
		if act.Distance.Kilometers() >= 100.0 {
			cen.Metric++
		}
		if act.Distance.Miles() >= 100.0 {
			cen.Imperial++
		}
		return true
	}, pass.Activities)
	return &Result{
		Count:         len(pass.Activities),
		Distance:      unit.Length(dst),
		ElevationGain: unit.Length(elv),
		MovingTime:    dur / (time.Second * 60 * 60),
		Centuries:     cen,
	}, nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "totals",
		Doc:  Doc,
		Run:  Run,
	}
}
