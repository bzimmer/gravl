package gravl

import (
	"errors"
	"fmt"
	"strconv"
	"time"

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
		client, err := noaa.NewClient(
			noaa.WithHTTPTracing(c.Bool("http-tracing")),
			noaa.WithTimeout(10*time.Second),
		)
		if err != nil {
			return err
		}
		var fcst *wx.Forecast
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
			fcst, e = client.Points.Forecast(c.Context, geom.Point[1], geom.Point[0])
			if e != nil {
				return e
			}
		case 2:
			lat, e := strconv.ParseFloat(args[0], 64)
			if e != nil {
				return e
			}
			lng, e := strconv.ParseFloat(args[1], 64)
			if e != nil {
				return e
			}
			fcst, e = client.Points.Forecast(c.Context, lat, lng)
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
			fcst, e = client.GridPoints.Forecast(c.Context, wfo, x, y)
			if e != nil {
				return e
			}
		default:
			return fmt.Errorf("only 2 or 3 arguments allowed [%v]", args)
		}
		err = encoder.Encode(fcst)
		if err != nil {
			return err
		}
		return nil
	},
}
