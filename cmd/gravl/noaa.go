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
			err := decoder.Decode(geom)
			if err != nil {
				return err
			}
			if !geom.IsPoint() {
				return errors.New("only Point geometries are supported")
			}
			fcst, err = client.Points.Forecast(c.Context, geom.Point[1], geom.Point[0])
		case 2:
			lat, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return err
			}
			lng, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return err
			}
			fcst, err = client.Points.Forecast(c.Context, lat, lng)
		case 3:
			wfo := args[0]
			x, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			y, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}
			fcst, err = client.GridPoints.Forecast(c.Context, wfo, x, y)
		default:
			return fmt.Errorf("only 2 or 3 arguments allowed [%v]", args)
		}
		if err != nil {
			return err
		}
		err = encoder.Encode(fcst)
		if err != nil {
			return err
		}
		return nil
	},
}
