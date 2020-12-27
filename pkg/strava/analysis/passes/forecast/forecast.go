package forecast

import (
	"context"
	"flag"

	"github.com/bzimmer/gravl/pkg/noaa"
	"github.com/bzimmer/gravl/pkg/strava/analysis"
	"github.com/twpayne/go-geom"
)

const Doc = ``

type Result struct {
	Activity *analysis.Activity `json:"activity"`
	Forecast *noaa.Forecast     `json:"forecast"`
}

type Forecast struct {
	Units  analysis.Units
	client *noaa.Client
}

func (a *Forecast) Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var res []*Result
	for _, act := range pass.Activities {
		coords := act.StartLatlng
		point := geom.NewPointFlat(geom.XY, []float64{coords[1], coords[0]})
		forecast, err := a.client.Points.Forecast(ctx, point)
		if err != nil {
			return nil, err
		}
		res = append(res, &Result{
			Activity: analysis.ToActivity(act, a.Units),
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
	c := &Forecast{
		Units:  analysis.Imperial,
		client: client,
	}
	fs := flag.NewFlagSet("forecast", flag.ExitOnError)
	fs.Var(&analysis.UnitsFlag{Units: &c.Units}, "units", "units to use")
	return &analysis.Analyzer{
		Name:  fs.Name(),
		Doc:   Doc,
		Flags: fs,
		Run:   c.Run,
	}
}
