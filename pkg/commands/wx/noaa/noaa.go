package noaa

import (
	"context"

	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/commands/encoding"
	wxcmd "github.com/bzimmer/gravl/pkg/commands/wx"
	"github.com/bzimmer/gravl/pkg/providers/wx/noaa"
)

func NewClient(c *cli.Context) (*noaa.Client, error) {
	return noaa.NewClient(noaa.WithHTTPTracing(c.Bool("http-tracing")))
}

func forecast(c *cli.Context) error {
	opts, err := wxcmd.Options(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	fcst, err := client.Points.Forecast(ctx, opts)
	if err != nil {
		return err
	}
	return encoding.Encode(fcst)
}

var Command = &cli.Command{
	Name:     "noaa",
	Category: "wx",
	Usage:    "Query NOAA for forecasts",
	Subcommands: []*cli.Command{
		{Name: "forecast", Action: forecast},
	},
}
