package gravl

import (
	"context"

	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/gnis"
)

var gnisCommand = &cli.Command{
	Name:     "gnis",
	Category: "geolocation",
	Usage:    "Query the GNIS database",
	Action: func(c *cli.Context) error {
		client, err := gnis.NewClient(
			gnis.WithHTTPTracing(c.Bool("http-tracing")),
		)
		if err != nil {
			return err
		}
		args := c.Args()
		for i := 0; i < args.Len(); i++ {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			features, err := client.GeoNames.Query(ctx, args.Get(i))
			if err != nil {
				return err
			}
			err = encoder.Encode(features)
			if err != nil {
				return err
			}
		}
		return nil
	},
}
