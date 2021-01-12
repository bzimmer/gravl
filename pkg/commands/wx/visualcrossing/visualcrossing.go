package visualcrossing

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/twpayne/go-geom"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/wx/visualcrossing"
)

func NewClient(c *cli.Context) (*visualcrossing.Client, error) {
	return visualcrossing.NewClient(
		visualcrossing.WithTokenCredentials(c.String("visualcrossing.access-token"), "", time.Time{}),
		visualcrossing.WithHTTPTracing(c.Bool("http-tracing")))
}

func forecast(c *cli.Context) error {
	opt := visualcrossing.ForecastOptions{
		Astronomy:      true,
		AggregateHours: c.Int("interval"),
		Units:          visualcrossing.UnitsMetric,
		AlertLevel:     visualcrossing.AlertLevelDetail,
	}
	switch c.Args().Len() {
	case 1:
		opt.Location = c.Args().Get(0)
	case 2:
		lng, err := strconv.ParseFloat(c.Args().Get(1), 64)
		if err != nil {
			return err
		}
		lat, err := strconv.ParseFloat(c.Args().Get(0), 64)
		if err != nil {
			return err
		}
		opt.Point = geom.NewPointFlat(geom.XY, []float64{lng, lat})
	default:
		return fmt.Errorf("only 1 or 2 arguments allowed [%v]", c.Args())
	}

	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	fcst, err := client.Forecast.Forecast(ctx, opt)
	if err != nil {
		return err
	}
	return encoding.Encode(fcst)
}

var forecastCommand = &cli.Command{
	Name: "forecast",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "interval",
			Aliases: []string{"i"},
			Value:   12,
			Usage:   "Forecast interval (eg 1, 12, 24)",
		},
	},
	Action: forecast,
}

var Command = &cli.Command{
	Name:     "visualcrossing",
	Aliases:  []string{"vc"},
	Category: "wx",
	Usage:    "Query VisualCrossing for forecasts",
	Flags:    AuthFlags,
	Subcommands: []*cli.Command{
		forecastCommand,
	},
}

var AuthFlags = []cli.Flag{
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "visualcrossing.access-token",
		Usage: "API key for Visual Crossing API",
	}),
}
