package gravl

import (
	"context"

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
			Usage: "API key for Visual Crossing API",
		}),
		&cli.IntFlag{
			Name:    "interval",
			Aliases: []string{"i"},
			Value:   12,
			Usage:   "Forecast interval (eg 1, 12, 24)",
		},
	},
	Action: func(c *cli.Context) error {
		client, err := visualcrossing.NewClient(
			visualcrossing.WithAPIKey(c.String("vc.api-key")),
			visualcrossing.WithHTTPTracing(c.Bool("http-tracing")),
		)
		if err != nil {
			return err
		}
		interval := c.Int("interval")
		for i := 0; i < c.Args().Len(); i++ {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			fcst, err := client.Forecast.Forecast(ctx,
				visualcrossing.WithLocation(c.Args().Get(i)),
				visualcrossing.WithAstronomy(true),
				visualcrossing.WithUnits(visualcrossing.UnitsMetric),
				visualcrossing.WithAggregateHours(interval),
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
