package rwgps

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/time/rate"

	api "github.com/bzimmer/activity"
	"github.com/bzimmer/activity/rwgps"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/activity"
)

const Provider = "rwgps"

var before sync.Once

func athlete(c *cli.Context) error {
	client := gravl.Runtime(c).RideWithGPS
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	user, err := client.Users.AuthenticatedUser(ctx)
	if err != nil {
		return err
	}
	log.Info().Int64("id", int64(user.ID)).Str("username", user.Name).Msg(c.Command.Name)
	gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, c.Command.Name}, 1)
	err = gravl.Runtime(c).Encoder.Encode(user)
	if err != nil {
		return err
	}
	return nil
}

func athleteCommand() *cli.Command {
	return &cli.Command{
		Name:    "athlete",
		Aliases: []string{"t"},
		Usage:   "Query for the authenticated athlete",
		Action:  athlete,
	}
}

func trips(c *cli.Context, kind string) error {
	client := gravl.Runtime(c).RideWithGPS
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	user, err := client.Users.AuthenticatedUser(ctx)
	if err != nil {
		return err
	}
	var metric string
	var trips []*rwgps.Trip
	switch kind {
	case "trips":
		metric = "activity"
		trips, err = client.Trips.Trips(ctx, user.ID, api.Pagination{Total: c.Int("count")})
	case "routes":
		metric = "route"
		trips, err = client.Trips.Routes(ctx, user.ID, api.Pagination{Total: c.Int("count")})
	default:
		trips, err = nil, fmt.Errorf("unknown type '%s'", kind)
	}
	if err != nil {
		return err
	}
	enc := gravl.Runtime(c).Encoder
	gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, c.Command.Name}, 1)
	for i, trip := range trips {
		gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, metric}, 1)
		log.Info().
			Time("date", trip.DepartedAt).
			Int64("id", trip.ID).
			Str("name", trip.Name).
			Msg(c.Command.Name)
		err = enc.Encode([]any{i, trip})
		if err != nil {
			return err
		}
	}
	return nil
}

func activitiesCommand() *cli.Command {
	return &cli.Command{
		Name:    "activities",
		Aliases: []string{"A"},
		Usage:   "Query activities for the authenticated athlete",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "count",
				Aliases: []string{"N"},
				Value:   0,
				Usage:   "The number of activities to query from RideWithGPS",
			},
		},
		Action: func(c *cli.Context) error { return trips(c, "trips") },
	}
}

func routesCommand() *cli.Command {
	return &cli.Command{
		Name:    "routes",
		Usage:   "Query routes for an athlete from RideWithGPS",
		Aliases: []string{"R"},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "count",
				Aliases: []string{"N"},
				Value:   0,
				Usage:   "The number of routes to query from RideWithGPS",
			},
		},
		Action: func(c *cli.Context) error { return trips(c, "routes") },
	}
}

func entity(c *cli.Context, f func(context.Context, int64) (any, error)) error {
	args := c.Args()
	enc := gravl.Runtime(c).Encoder
	for i := 0; i < args.Len(); i++ {
		err := func() error {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			x, err := strconv.ParseInt(args.Get(i), 0, 0)
			if err != nil {
				return err
			}
			v, err := f(ctx, x)
			if err != nil {
				return err
			}
			gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, c.Command.Name}, 1)
			if err := enc.Encode(v); err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}

func activityCommand() *cli.Command { //nolint:dupl
	return &cli.Command{
		Name:    "activity",
		Aliases: []string{"a"},
		Usage:   "Query an activity from RideWithGPS",
		Action: func(c *cli.Context) error {
			client := gravl.Runtime(c).RideWithGPS
			return entity(c, func(ctx context.Context, id int64) (any, error) {
				trip, err := client.Trips.Trip(ctx, id)
				if err != nil {
					return nil, err
				}
				log.Info().Int64("id", trip.ID).Str("name", trip.Name).Msg(c.Command.Name)
				return trip, nil
			})
		},
	}
}

func routeCommand() *cli.Command { //nolint:dupl
	return &cli.Command{
		Name:    "route",
		Aliases: []string{"r"},
		Usage:   "Query a route from RideWithGPS",
		Action: func(c *cli.Context) error {
			client := gravl.Runtime(c).RideWithGPS
			return entity(c, func(ctx context.Context, id int64) (any, error) {
				route, err := client.Trips.Route(ctx, id)
				if err != nil {
					return nil, err
				}
				log.Info().Int64("id", route.ID).Str("name", route.Name).Msg(c.Command.Name)
				return route, nil
			})
		},
	}
}

func Before(c *cli.Context) error {
	var err error
	before.Do(func() {
		var client *rwgps.Client
		client, err = rwgps.NewClient(
			rwgps.WithClientCredentials(c.String("rwgps-client-id"), ""),
			rwgps.WithTokenCredentials(c.String("rwgps-access-token"), "", time.Time{}),
			rwgps.WithHTTPTracing(c.Bool("http-tracing")),
			rwgps.WithRateLimiter(rate.NewLimiter(
				rate.Every(c.Duration("rate-limit")), c.Int("rate-burst"))))
		if err != nil {
			return
		}
		gravl.Runtime(c).RideWithGPS = client
		gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, "client", "created"}, 1)
		log.Info().Msg("created rwgps client")
	})
	return err
}

func Command() *cli.Command {
	return &cli.Command{
		Name:        "rwgps",
		Category:    "activity",
		Usage:       "Query RideWithGPS for rides and routes",
		Description: "Operations supported by the RideWithGPS API",
		Flags:       append(AuthFlags(), activity.RateLimitFlags()...),
		Before:      Before,
		Subcommands: []*cli.Command{
			activitiesCommand(),
			activityCommand(),
			athleteCommand(),
			routeCommand(),
			routesCommand(),
		},
	}
}

func AuthFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "rwgps-client-id",
			Usage:   "RideWithGPS client id",
			EnvVars: []string{"RWGPS_CLIENT_ID"},
		},
		&cli.StringFlag{
			Name:    "rwgps-access-token",
			Usage:   "RideWithGPS access token",
			EnvVars: []string{"RWGPS_ACCESS_TOKEN"},
		},
	}
}
