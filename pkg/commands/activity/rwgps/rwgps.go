package rwgps

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/activity"
	"github.com/bzimmer/gravl/pkg/providers/activity/rwgps"
)

func NewClient(c *cli.Context) (*rwgps.Client, error) {
	return rwgps.NewClient(
		rwgps.WithClientCredentials(c.String("rwgps.client-id"), ""),
		rwgps.WithTokenCredentials(c.String("rwgps.access-token"), "", time.Time{}),
		rwgps.WithHTTPTracing(c.Bool("http-tracing")),
	)
}

func athlete(c *cli.Context) error {
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	user, err := client.Users.AuthenticatedUser(ctx)
	if err != nil {
		return err
	}
	err = encoding.Encode(user)
	if err != nil {
		return err
	}
	return nil
}

var athleteCommand = &cli.Command{
	Name:    "athlete",
	Aliases: []string{"t"},
	Usage:   "Query for the authenticated athlete",
	Action:  athlete,
}

func trips(c *cli.Context, kind string) error {
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	user, err := client.Users.AuthenticatedUser(ctx)
	if err != nil {
		return err
	}
	var trips []*rwgps.Trip
	switch kind {
	case "trips":
		trips, err = client.Trips.Trips(ctx, user.ID, activity.Pagination{Total: c.Int("count")})
	case "routes":
		trips, err = client.Trips.Routes(ctx, user.ID, activity.Pagination{Total: c.Int("count")})
	default:
		return fmt.Errorf("unknown kind '%s'", kind)
	}
	if err != nil {
		return err
	}
	for _, trip := range trips {
		err = encoding.Encode(trip)
		if err != nil {
			return err
		}
	}
	return nil
}

var activitiesCommand = &cli.Command{
	Name:    "activities",
	Aliases: []string{"A"},
	Usage:   "Query activities for the authenticated athlete",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "count",
			Aliases: []string{"N"},
			Value:   0,
			Usage:   "Count",
		},
	},
	Action: func(c *cli.Context) error { return trips(c, "trips") },
}

func entity(c *cli.Context, f func(context.Context, *rwgps.Client, int64) (interface{}, error)) error {
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	for i := 0; i < c.Args().Len(); i++ {
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		x, err := strconv.ParseInt(c.Args().Get(i), 0, 0)
		if err != nil {
			return err
		}
		v, err := f(ctx, client, x)
		if err != nil {
			return err
		}
		if err = encoding.Encode(v); err != nil {
			return err
		}
	}
	return nil
}

var activityCommand = &cli.Command{
	Name:    "activity",
	Aliases: []string{"a"},
	Usage:   "Query an activity from RideWithGPS",
	Action: func(c *cli.Context) error {
		return entity(c, func(ctx context.Context, client *rwgps.Client, id int64) (interface{}, error) {
			return client.Trips.Trip(ctx, id)
		})
	},
}

var routeCommand = &cli.Command{
	Name:    "route",
	Aliases: []string{"r"},
	Usage:   "Query a route from RideWithGPS",
	Action: func(c *cli.Context) error {
		return entity(c, func(ctx context.Context, client *rwgps.Client, id int64) (interface{}, error) {
			return client.Trips.Route(ctx, id)
		})
	},
}

var routesCommand = &cli.Command{
	Name:    "routes",
	Usage:   "Query routes for an athlete from RideWithGPS",
	Aliases: []string{"R"},
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "count",
			Aliases: []string{"N"},
			Value:   0,
			Usage:   "Count",
		},
	},
	Action: func(c *cli.Context) error { return trips(c, "routes") },
}

var Command = &cli.Command{
	Name:     "rwgps",
	Category: "activity",
	Usage:    "Query RideWithGPS for rides and routes",
	Flags:    AuthFlags,
	Subcommands: []*cli.Command{
		activitiesCommand,
		activityCommand,
		athleteCommand,
		routeCommand,
		routesCommand,
	},
}

var AuthFlags = []cli.Flag{
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "rwgps.client-id",
		Value: "",
		Usage: "Client ID for RideWithGPS API",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "rwgps.access-token",
		Value: "",
		Usage: "Access token for RideWithGPS API",
	}),
}
