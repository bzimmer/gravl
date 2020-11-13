package gravl

import (
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
		var (
			err error
			rte *route.Route
		)
		client, err := rwgps.NewClient(
			rwgps.WithAPIKey(c.String("rwgps.api-key")),
			rwgps.WithAuthToken(c.String("rwgps.auth-token")),
			rwgps.WithHTTPTracing(c.Bool("http-tracing")),
		)
		if err != nil {
			return err
		}
		if c.Bool("athlete") {
			user, err := client.Users.AuthenticatedUser(c.Context)
			if err != nil {
				return err
			}
			err = encoder.Encode(user)
			if err != nil {
				return err
			}
			return nil
		}
		for i := 0; i < c.Args().Len(); i++ {
			x, err := strconv.ParseInt(c.Args().Get(i), 0, 0)
			if err != nil {
				return err
			}
			if c.Bool("trip") {
				rte, err = client.Trips.Trip(c.Context, x)
			} else {
				rte, err = client.Trips.Route(c.Context, x)
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
