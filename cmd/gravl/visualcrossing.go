package gravl

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/common/wx"
	"github.com/bzimmer/gravl/pkg/visualcrossing"
)

var visualcrossingCommand = &cli.Command{
	Name:     "visualcrossing",
	Aliases:  []string{"vc"},
	Category: "forecast",
	Usage:    "Query VisualCrossing for forecasts",
	Flags: []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "visualcrossing.access-token",
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
			visualcrossing.WithTokenCredentials(c.String("visualcrossing.access-token"), "", time.Time{}),
			visualcrossing.WithHTTPTracing(c.Bool("http-tracing")))
		if err != nil {
			return err
		}

		var fcst wx.Forecastable
		var opt = visualcrossing.ForecastOptions{
			Astronomy:      true,
			AggregateHours: c.Int("interval"),
			Units:          visualcrossing.UnitsMetric,
			AlertLevel:     visualcrossing.AlertLevelDetail,
		}

		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()

		switch c.Args().Len() {
		case 1:
			opt.Location = c.Args().Get(0)
		case 2:
			lng, e := strconv.ParseFloat(c.Args().Get(0), 64)
			if e != nil {
				return e
			}
			lat, e := strconv.ParseFloat(c.Args().Get(1), 64)
			if e != nil {
				return e
			}
			opt.Coordinates = visualcrossing.Coordinates{
				Latitude:  lat,
				Longitude: lng,
			}
		default:
			return fmt.Errorf("only 1 or 2 arguments allowed [%v]", c.Args())
		}
		fcst, err = client.Forecast.Forecast(ctx, opt)
		if err != nil {
			return err
		}
		return encoder.Encode(fcst)
	},
}
