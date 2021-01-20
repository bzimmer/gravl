package totals

import (
	"github.com/martinlindhe/unit"

	"github.com/bzimmer/gravl/pkg/analysis"
)

const doc = `totals returns the number of centuries (100 mi or 100 km).`

type Centuries struct {
	Metric   int `json:"metric"`
	Imperial int `json:"imperial"`
}

type Result struct {
	Count      int           `json:"count"`
	Distance   float64       `json:"distance"`
	Elevation  float64       `json:"elevation"`
	MovingTime unit.Duration `json:"movingtime"`
	Centuries  Centuries     `json:"centuries"`
}

func run(ctx *analysis.Context, pass *analysis.Pass) (interface{}, error) {
	var cen Centuries
	var dst, elv float64
	var dur unit.Duration
	for i := 0; i < len(pass.Activities); i++ {
		act := pass.Activities[i]
		dur = dur + act.MovingTime
		switch ctx.Units {
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
	}
	return &Result{
		Count:      len(pass.Activities),
		Distance:   dst,
		Elevation:  elv,
		MovingTime: dur,
		Centuries:  cen,
	}, nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "totals",
		Doc:  doc,
		Run:  run,
	}
}
