package gravl

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	geojson "github.com/paulmach/go.geojson"
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
		var fcst wx.Forecastable
		args := c.Args().Slice()
		switch len(args) {
		case 0:
			geom := &geojson.Geometry{}
			e := decoder.Decode(geom)
			if e != nil {
				return e
			}
			if !geom.IsPoint() {
				return errors.New("only Point geometries are supported")
			}
			fcst, e = client.Points.Forecast(ctx, geom.Point[1], geom.Point[0])
			if e != nil {
				return e
			}
		case 2:
			lng, e := strconv.ParseFloat(args[0], 64)
			if e != nil {
				return e
			}
			lat, e := strconv.ParseFloat(args[1], 64)
			if e != nil {
				return e
			}
			fcst, e = client.Points.Forecast(ctx, lat, lng)
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
		if err = encoder.Encode(f); err != nil {
			return err
		}
		return nil
	},
}
