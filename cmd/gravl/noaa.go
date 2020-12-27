package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/twpayne/go-geom"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/common/wx"
	"github.com/bzimmer/gravl/pkg/noaa"
)

var noaaCommand = &cli.Command{
	Name:     "noaa",
	Category: "forecast",
	Usage:    "Query NOAA for forecasts",
	Action: func(c *cli.Context) error {
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		client, err := noaa.NewClient(
			noaa.WithHTTPTracing(c.Bool("http-tracing")),
		)
		if err != nil {
			return err
		}
		var fcst wx.Forecaster
		args := c.Args().Slice()
		switch len(args) {
		case 2:
			lng, e := strconv.ParseFloat(args[0], 64)
			if e != nil {
				return e
			}
			lat, e := strconv.ParseFloat(args[1], 64)
			if e != nil {
				return e
			}
			point := geom.NewPointFlat(geom.XY, []float64{lng, lat})
			fcst, e = client.Points.Forecast(ctx, point)
			if e != nil {
				return e
			}
		case 3:
			wfo := args[0]
			x, e := strconv.Atoi(args[1])
			if e != nil {
				return e
			}
			y, e := strconv.Atoi(args[2])
			if e != nil {
				return e
			}
			fcst, e = client.GridPoints.Forecast(ctx, wfo, x, y)
			if e != nil {
				return e
			}
		default:
			return fmt.Errorf("only 2 or 3 arguments allowed [%v]", args)
		}
		f, err := fcst.Forecast()
		if err != nil {
			return err
		}
		return encoder.Encode(f)
	},
}
