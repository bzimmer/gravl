package wx

import (
	"fmt"
	"strconv"

	"github.com/twpayne/go-geom"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/providers/wx"
)

// Options parses the cli flags and args to create forecast options
func Options(c *cli.Context) (wx.ForecastOptions, error) {
	opts := wx.ForecastOptions{
		Units:          wx.Metric,
		AggregateHours: c.Int("interval"),
	}
	switch c.NArg() {
	case 1:
		opts.Location = c.Args().Get(0)
	case 2:
		lng, err := strconv.ParseFloat(c.Args().Get(1), 64)
		if err != nil {
			return wx.ForecastOptions{}, err
		}
		lat, err := strconv.ParseFloat(c.Args().Get(0), 64)
		if err != nil {
			return wx.ForecastOptions{}, err
		}
		opts.Point = geom.NewPointFlat(geom.XY, []float64{lng, lat})
	default:
		return wx.ForecastOptions{}, fmt.Errorf("expected 1 or 2 arguments, found [%d]", c.NArg())
	}
	return opts, nil
}
