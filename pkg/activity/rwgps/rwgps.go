package rwgps

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
	"golang.org/x/time/rate"

	"github.com/bzimmer/activity"
	"github.com/bzimmer/activity/rwgps"
	"github.com/bzimmer/gravl/pkg"
	actcmd "github.com/bzimmer/gravl/pkg/activity"
)

const provider = "rwgps"

func athlete(c *cli.Context) error {
	client := pkg.Runtime(c).RideWithGPS
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	user, err := client.Users.AuthenticatedUser(ctx)
	if err != nil {
		return err
	}
	pkg.Runtime(c).Metrics.IncrCounter([]string{provider, c.Command.Name}, 1)
	err = pkg.Runtime(c).Encoder.Encode(user)
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
	client := pkg.Runtime(c).RideWithGPS
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
		return fmt.Errorf("unknown type '%s'", kind)
	}
	if err != nil {
		return err
	}
	enc := pkg.Runtime(c).Encoder
	for _, trip := range trips {
		err = enc.Encode(trip)
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
			Usage:   "The number of activities to query from RideWithGPS (the number returned will be <= N)",
		},
	},
	Action: func(c *cli.Context) error { return trips(c, "trips") },
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
			Usage:   "The number of routes to query from RideWithGPS (the number returned will be <= N)",
		},
	},
	Action: func(c *cli.Context) error { return trips(c, "routes") },
}

func entity(c *cli.Context, f func(context.Context, *rwgps.Client, int64) (interface{}, error)) error {
	enc := pkg.Runtime(c).Encoder
	client := pkg.Runtime(c).RideWithGPS
	args := c.Args()
	for i := 0; i < args.Len(); i++ {
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		x, err := strconv.ParseInt(args.Get(i), 0, 0)
		if err != nil {
			return err
		}
		v, err := f(ctx, client, x)
		if err != nil {
			return err
		}
		if err := enc.Encode(v); err != nil {
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

func Before(c *cli.Context) error {
	client, err := rwgps.NewClient(
		rwgps.WithClientCredentials(c.String("rwgps-client-id"), ""),
		rwgps.WithTokenCredentials(c.String("rwgps-access-token"), "", time.Time{}),
		rwgps.WithHTTPTracing(c.Bool("http-tracing")),
		rwgps.WithRateLimiter(rate.NewLimiter(
			rate.Every(c.Duration("rate-limit")), c.Int("rate-burst"))))
	if err != nil {
		return err
	}
	pkg.Runtime(c).RideWithGPS = client
	return nil
}

func Command() *cli.Command {
	return &cli.Command{
		Name:        "rwgps",
		Category:    "activity",
		Usage:       "Query RideWithGPS for rides and routes",
		Description: "Operations supported by the RideWithGPS API",
		Flags:       append(AuthFlags, actcmd.RateLimitFlags...),
		Before:      Before,
		Subcommands: []*cli.Command{
			activitiesCommand,
			activityCommand,
			athleteCommand,
			routeCommand,
			routesCommand,
			actcmd.UploadCommand(func(c *cli.Context) (activity.Uploader, error) {
				return pkg.Runtime(c).RideWithGPS.Uploader(), nil
			}),
		},
	}
}

var AuthFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "rwgps-client-id",
		Required: true,
		Usage:    "rwgps client id",
		EnvVars:  []string{"RWGPS_CLIENT_ID"},
	},
	&cli.StringFlag{
		Name:     "rwgps-access-token",
		Required: true,
		Usage:    "rwgps access token",
		EnvVars:  []string{"RWGPS_ACCESS_TOKEN"},
	},
}
