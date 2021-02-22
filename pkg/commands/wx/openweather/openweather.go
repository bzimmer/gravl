package openweather

import (
	"context"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/commands/wx"
	"github.com/bzimmer/gravl/pkg/providers/wx/openweather"
)

func NewClient(c *cli.Context) (*openweather.Client, error) {
	return openweather.NewClient(
		openweather.WithTokenCredentials(c.String("openweather.access-token"), "", time.Time{}),
		openweather.WithHTTPTracing(c.Bool("http-tracing")))
}

func forecast(c *cli.Context) error {
	opts, err := wx.Options(c)
	if err != nil {
		return err
	}
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	fcst, err := client.Forecast.Forecast(ctx, opts)
	if err != nil {
		return err
	}
	return encoding.Encode(fcst)
}

var Command = &cli.Command{
	Name:     "openweather",
	Aliases:  []string{"ow"},
	Category: "wx",
	Usage:    "Query OpenWeather for forecasts",
	Flags:    AuthFlags,
	Subcommands: []*cli.Command{
		{
			Name:      "forecast",
			Usage:     "Query OpenWeather for a forecast",
			ArgsUsage: "[--] LATITUTE LONGITUDE",
			Action:    forecast,
		},
	},
}

var AuthFlags = []cli.Flag{
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "openweather.access-token",
		Usage: "API key for OpenWeather API",
	}),
}
