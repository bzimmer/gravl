package srtm

import (
	"context"
	"fmt"
	"path"
	"strconv"

	"github.com/adrg/xdg"
	"github.com/twpayne/go-geom"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/geo/srtm"
)

var Command = &cli.Command{
	Name:     "srtm",
	Category: "geo",
	Usage:    "Query the SRTM database for elevation data",
	Action: func(c *cli.Context) error {
		var err error
		var longitude, latitude float64
		switch c.Args().Len() {
		case 2:
			longitude, err = strconv.ParseFloat(c.Args().Get(1), 64)
			if err != nil {
				return err
			}
			latitude, err = strconv.ParseFloat(c.Args().Get(0), 64)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("expected <lat> <lng>, found [%v]", c.Args())
		}
		client, err := srtm.NewClient(
			srtm.WithHTTPTracing(c.Bool("http-tracing")),
			srtm.WithStorageLocation(path.Join(xdg.CacheHome, pkg.PackageName, c.Command.Name)))
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		pt := geom.NewPointFlat(geom.XY, []float64{longitude, latitude})
		elevation, err := client.Elevation.Elevation(ctx, pt)
		if err != nil {
			return err
		}
		return encoding.For(c).Encode(elevation)
	},
}
