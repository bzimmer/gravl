package forecast

import (
	"context"

	"github.com/twpayne/go-geom"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/providers/wx/noaa"
)

const doc = ``

type Result struct {
	Activity *analysis.Activity `json:"activity"`
	Forecast *noaa.Forecast     `json:"forecast"`
}

type forecast struct {
	client *noaa.Client
}

func (a *forecast) run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var res []*Result
	for _, act := range pass.Activities {
		coords := act.StartLatlng
		point := geom.NewPointFlat(geom.XY, []float64{coords[1], coords[0]})
		forecast, err := a.client.Points.Forecast(ctx, point)
		if err != nil {
			return nil, err
		}
		res = append(res, &Result{
			Activity: analysis.ToActivity(act, pass.Units),
			Forecast: forecast,
		})
	}
	return res, nil
}

func New() *analysis.Analyzer {
	client, err := noaa.NewClient()
	if err != nil {
		panic(err)
	}
	f := &forecast{client: client}
	return &analysis.Analyzer{
		Name: "forecast",
		Doc:  doc,
		Run:  f.run,
	}
}
