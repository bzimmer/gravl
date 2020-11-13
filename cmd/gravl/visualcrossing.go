package gravl

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/visualcrossing"
)

var visualcrossingCommand = &cli.Command{
	Name:     "vc",
	Category: "forecast",
	Usage:    "Query VisualCrossing for forecasts",
	Flags: []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "vc.api-key",
			Usage: "API key for VC API",
		}),
	},
	Action: func(c *cli.Context) error {
		client, err := visualcrossing.NewClient(
			visualcrossing.WithAPIKey(c.String("vc.api-key")),
			visualcrossing.WithHTTPTracing(c.Bool("http-tracing")),
		)
		if err != nil {
			return err
		}
		for i := 0; i < c.Args().Len(); i++ {
			fcst, err := client.Forecast.Forecast(
				c.Context,
				visualcrossing.WithLocation(c.Args().Get(i)),
				visualcrossing.WithAstronomy(true),
				visualcrossing.WithUnits(visualcrossing.UnitsMetric),
				visualcrossing.WithAggregateHours(1),
				visualcrossing.WithAlerts(visualcrossing.AlertLevelDetail))
			if err != nil {
				return err
			}
			err = encoder.Encode(fcst)
			if err != nil {
				return err
			}
		}
		return nil
	},
}
