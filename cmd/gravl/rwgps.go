package gravl

import (
	"context"
	"strconv"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/common/route"
	"github.com/bzimmer/gravl/pkg/rwgps"
)

var rwgpsCommand = &cli.Command{
	Name:     "rwgps",
	Category: "route",
	Usage:    "Query Ride with GPS for rides and routes",
	Flags: []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "rwgps.api-key",
			Value: "",
			Usage: "Access key for RWGPS API",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "rwgps.auth-token",
			Value: "",
			Usage: "Auth token for RWGPS API",
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
			rwgps.WithAPIKey(c.String("rwgps.api-key")),
			rwgps.WithAuthToken(c.String("rwgps.auth-token")),
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

		var rte *route.Route
		for i := 0; i < c.Args().Len(); i++ {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			x, err := strconv.ParseInt(c.Args().Get(i), 0, 0)
			if err != nil {
				return err
			}
			if c.Bool("trip") {
				rte, err = client.Trips.Trip(ctx, x)
			} else {
				rte, err = client.Trips.Route(ctx, x)
			}
			if err != nil {
				return err
			}
			err = encoder.Encode(rte)
			if err != nil {
				return err
			}
		}
		return nil
	},
}
