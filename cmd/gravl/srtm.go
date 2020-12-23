package main

import (
	"context"
	"fmt"
	"path"
	"strconv"

	"github.com/adrg/xdg"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/srtm"
)

const (
	srtmCache = "srtm"
)

var srtmCommand = &cli.Command{
	Name:     "srtm",
	Category: "geolocation",
	Usage:    "Query the SRTM database for elevation data",
	Action: func(c *cli.Context) error {
		var (
			err                 error
			longitude, latitude float64
		)
		switch c.Args().Len() {
		case 0:
			// Barlow Pass
			longitude, latitude = -121.4440005, 48.0264959
		case 2:
			longitude, err = strconv.ParseFloat(c.Args().Get(0), 64)
			if err != nil {
				return err
			}
			latitude, err = strconv.ParseFloat(c.Args().Get(1), 64)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("only 0 or 2 arguments allowed [%v]", c.Args())
		}

		client, err := srtm.NewClient(
			srtm.WithHTTPTracing(c.Bool("http-tracing")),
			srtm.WithStorageLocation(
				path.Join(xdg.CacheHome, pkg.PackageName, srtmCache)),
		)
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		elevation, err := client.Elevation.Elevation(ctx, longitude, latitude)
		if err != nil {
			return err
		}
		return encoder.Encode(elevation)
	},
}
