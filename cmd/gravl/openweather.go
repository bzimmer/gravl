package gravl

import (
	"context"
	"fmt"
	"strconv"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/openweather"
)

var openweatherCommand = &cli.Command{
	Name:     "ow",
	Category: "forecast",
	Usage:    "Query OpenWeather for forecasts",
	Flags: []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "ow.api-key",
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
		client, err := openweather.NewClient(
			openweather.WithAPIKey(c.String("ow.api-key")),
			openweather.WithHTTPTracing(c.Bool("http-tracing")),
		)
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()

		lng, err := strconv.ParseFloat(c.Args().Get(0), 64)
		if err != nil {
			return err
		}
		lat, err := strconv.ParseFloat(c.Args().Get(1), 64)
		if err != nil {
			return err
		}
		fcst, err := client.Forecast.Forecast(ctx, openweather.WithLocation(lng, lat))
		if err != nil {
			return err
		}
		return encoder.Encode(fcst)
	},
}
