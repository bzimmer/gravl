package gravl

import (
	"context"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/common/geo"
	"github.com/bzimmer/gravl/pkg/rwgps"
)

var rwgpsCommand = &cli.Command{
	Name:     "rwgps",
	Category: "route",
	Usage:    "Query Ride with GPS for rides and routes",
	Flags: []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "rwgps.client-id",
			Value: "",
			Usage: "Client ID for RWGPS API",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "rwgps.access-token",
			Value: "",
			Usage: "Access token for RWGPS API",
		}),
		&cli.BoolFlag{
			Name:    "trip",
			Aliases: []string{"t"},
			Value:   false,
			Usage:   "Trip",
		},
		&cli.BoolFlag{
			Name:    "route",
			Aliases: []string{"r"},
			Value:   false,
			Usage:   "Route",
		},
		&cli.BoolFlag{
			Name:    "athlete",
			Aliases: []string{"a"},
			Value:   false,
			Usage:   "Athlete",
		},
	},
	Action: func(c *cli.Context) error {
		client, err := rwgps.NewClient(
			rwgps.WithClientCredentials(c.String("rwgps.client-id"), ""),
			rwgps.WithTokenCredentials(c.String("rwgps.access-token"), "", time.Time{}),
			rwgps.WithHTTPTracing(c.Bool("http-tracing")),
		)
		if err != nil {
			return err
		}

		if c.Bool("athlete") {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			user, err := client.Users.AuthenticatedUser(ctx)
			if err != nil {
				return err
			}
			err = encoder.Encode(user)
			if err != nil {
				return err
			}
			return nil
		}

		var tck geo.Tracker
		for i := 0; i < c.Args().Len(); i++ {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			x, err := strconv.ParseInt(c.Args().Get(i), 0, 0)
			if err != nil {
				return err
			}
			if c.Bool("trip") {
				tck, err = client.Trips.Trip(ctx, x)
			} else {
				tck, err = client.Trips.Route(ctx, x)
			}
			if err != nil {
				return err
			}
			t, err := tck.Track()
			if err != nil {
				return err
			}
			if err = encoder.Encode(t); err != nil {
				return err
			}
		}
		return nil
	},
}
