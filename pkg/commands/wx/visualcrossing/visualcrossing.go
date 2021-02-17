package visualcrossing

import (
	"context"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/commands/wx"
	"github.com/bzimmer/gravl/pkg/providers/wx/visualcrossing"
)

func NewClient(c *cli.Context) (*visualcrossing.Client, error) {
	return visualcrossing.NewClient(
		visualcrossing.WithTokenCredentials(c.String("visualcrossing.access-token"), "", time.Time{}),
		visualcrossing.WithHTTPTracing(c.Bool("http-tracing")))
}

func forecast(c *cli.Context) error {
	opts, err := wx.Options(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	fcst, err := client.Forecast.Forecast(ctx, opts)
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
	Name:        "visualcrossing",
	Aliases:     []string{"vc"},
	Category:    "wx",
	Usage:       "Query VisualCrossing for forecasts",
	Flags:       AuthFlags,
	Subcommands: []*cli.Command{forecastCommand},
}

var AuthFlags = []cli.Flag{
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "visualcrossing.access-token",
		Usage: "API key for Visual Crossing API",
	}),
}
