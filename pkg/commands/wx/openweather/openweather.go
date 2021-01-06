package openweather

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/twpayne/go-geom"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/wx/openweather"
)

var Command = &cli.Command{
	Name:     "openweather",
	Aliases:  []string{"ow"},
	Category: "wx",
	Usage:    "Query OpenWeather for forecasts",
	Flags:    AuthFlags,
	Before: func(c *cli.Context) error {
		if c.Args().Len() != 2 {
			return fmt.Errorf("invalid number of arguments, expected 2, got {%d}", c.Args().Len())
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		var err error
		var longitude, latitude float64
		longitude, err = strconv.ParseFloat(c.Args().Get(0), 64)
		if err != nil {
			return err
		}
		latitude, err = strconv.ParseFloat(c.Args().Get(1), 64)
		if err != nil {
			return err
		}
		client, err := openweather.NewClient(
			openweather.WithTokenCredentials(c.String("openweather.access-token"), "", time.Time{}),
			openweather.WithHTTPTracing(c.Bool("http-tracing")))
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		fcst, err := client.Forecast.Forecast(ctx,
			openweather.ForecastOptions{
				Units: openweather.UnitsMetric,
				Point: geom.NewPointFlat(geom.XY, []float64{longitude, latitude})})
		if err != nil {
			return err
		}
		return encoding.Encode(fcst)
	},
}

var AuthFlags = []cli.Flag{
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "openweather.access-token",
		Usage: "API key for OpenWeather API",
	}),
}