package gravl

import (
	"context"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/common"
	"github.com/bzimmer/gravl/pkg/rwgps"
)

var rwgpsCommand = &cli.Command{
	Name:     "rwgps",
	Category: "route",
	Usage:    "Query Ride with GPS for rides and routes",
	Flags: merge(
		rwgpsAuthFlags,
		[]cli.Flag{
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
			&cli.BoolFlag{
				Name:    "activities",
				Aliases: []string{"A"},
				Value:   false,
				Usage:   "Activities",
			},
			// &cli.BoolFlag{
			// 	Name:    "routes",
			// 	Aliases: []string{"R"},
			// 	Value:   false,
			// 	Usage:   "Routes",
			// },
			&cli.IntFlag{
				Name:    "count",
				Aliases: []string{"N"},
				Value:   10,
				Usage:   "Count",
			},
		}),
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

		if c.Bool("activities") {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			user, err := client.Users.AuthenticatedUser(ctx)
			if err != nil {
				return err
			}
			trips, err := client.Trips.Trips(ctx, user.ID, common.Pagination{Total: c.Int("count")})
			if err != nil {
				return err
			}
			for _, trip := range trips {
				err = encoder.Encode(trip)
				if err != nil {
					return err
				}
			}
			return nil
		}

		for i := 0; i < c.Args().Len(); i++ {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			x, err := strconv.ParseInt(c.Args().Get(i), 0, 0)
			if err != nil {
				return err
			}
			var t interface{}
			if c.Bool("trip") {
				t, err = client.Trips.Trip(ctx, x)
			} else {
				t, err = client.Trips.Route(ctx, x)
			}
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

var rwgpsAuthFlags = []cli.Flag{
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
}
