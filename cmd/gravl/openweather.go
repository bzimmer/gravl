package gravl

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/openweather"
)

var openweatherCommand = &cli.Command{
	Name:     "openweather",
	Aliases:  []string{"ow"},
	Category: "forecast",
	Usage:    "Query OpenWeather for forecasts",
	Flags: []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "openweather.access-token",
			Usage: "API key for OpenWeather API",
		}),
	},
	Before: func(c *cli.Context) error {
		if c.Args().Len() != 2 {
			return fmt.Errorf("invalid number of arguments, expected 2, got {%d}", c.Args().Len())
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		var (
			err                 error
			longitude, latitude float64
		)
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
				Coordinates: openweather.Coordinates{
					Latitude:  latitude,
					Longitude: longitude}})
		if err != nil {
			return err
		}
		return encoder.Encode(fcst)
	},
}
