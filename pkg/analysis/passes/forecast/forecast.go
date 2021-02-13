package forecast

import (
	"github.com/twpayne/go-geom"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
	"github.com/bzimmer/gravl/pkg/providers/wx"
	"github.com/bzimmer/gravl/pkg/providers/wx/noaa"
)

const doc = `forecast the weather for an activity`

type Result struct {
	Activity *analysis.Activity `json:"activity"`
	Forecast *noaa.Forecast     `json:"forecast"`
}

type forecast struct {
	client *noaa.Client
}

func (a *forecast) run(ctx *analysis.Context, pass []*strava.Activity) (interface{}, error) {
	var res []*Result
	for _, act := range pass {
		coords := act.StartLatlng
		opts := wx.ForecastOptions{
			Point: geom.NewPointFlat(geom.XY, []float64{coords[1], coords[0]}),
		}
		forecast, err := a.client.Points.Forecast(ctx, opts)
		if err != nil {
			return nil, err
		}
		res = append(res, &Result{
			Activity: analysis.ToActivity(act, ctx.Units),
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
