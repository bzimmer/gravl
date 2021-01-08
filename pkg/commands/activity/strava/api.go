package strava

import (
	"context"
	"strconv"
	"time"

	"github.com/bzimmer/gravl/pkg/activity"
	"github.com/bzimmer/gravl/pkg/activity/strava"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/time/rate"
)

type entityFunc func(context.Context, *strava.Client, int64) (interface{}, error)

func NewAPIClient(c *cli.Context) (*strava.Client, error) {
	return strava.NewClient(
		strava.WithTokenCredentials(
			c.String("strava.access-token"), c.String("strava.refresh-token"), time.Now().Add(-1*time.Minute)),
		strava.WithClientCredentials(c.String("strava.client-id"), c.String("strava.client-secret")),
		strava.WithAutoRefresh(c.Context),
		strava.WithHTTPTracing(c.Bool("http-tracing")),
		strava.WithRateLimiter(rate.NewLimiter(rate.Every(1500*time.Millisecond), 25)))
}

func athlete(c *cli.Context) error {
	client, err := NewAPIClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	athlete, err := client.Athlete.Athlete(ctx)
	if err != nil {
		return err
	}
	return encoding.Encode(athlete)
}

var athleteCommand = &cli.Command{
	Name:    "athlete",
	Usage:   "Query an athlete from Strava",
	Aliases: []string{"t"},
	Action:  athlete,
}

func refresh(c *cli.Context) error {
	client, err := NewAPIClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	tokens, err := client.Auth.Refresh(ctx)
	if err != nil {
		return err
	}
	return encoding.Encode(tokens)
}

var refreshCommand = &cli.Command{
	Name:   "refresh",
	Usage:  "Acquire a new refresh token",
	Action: refresh,
}

func activities(c *cli.Context) error {
	client, err := NewAPIClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	defer func(start time.Time) {
		log.Debug().
			Dur("elapsed", time.Since(start)).
			Msg("activities")
	}(time.Now())
	activities, err := client.Activity.Activities(ctx, activity.Pagination{Total: c.Int("count")})
	if err != nil {
		return err
	}
	for _, act := range activities {
		err = encoding.Encode(act)
		if err != nil {
			return err
		}
	}
	return nil
}

var activitiesCommand = &cli.Command{
	Name:    "activities",
	Usage:   "Query activities for an athlete from Strava",
	Aliases: []string{"A"},
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "count",
			Aliases: []string{"N"},
			Value:   0,
			Usage:   "Count",
		},
	},
	Action: activities,
}

func routes(c *cli.Context) error {
	client, err := NewAPIClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	defer func(start time.Time) {
		log.Debug().
			Dur("elapsed", time.Since(start)).
			Msg("activities")
	}(time.Now())
	athlete, err := client.Athlete.Athlete(ctx)
	if err != nil {
		return err
	}
	routes, err := client.Route.Routes(ctx, athlete.ID, activity.Pagination{Total: c.Int("count")})
	if err != nil {
		return err
	}
	for _, route := range routes {
		err = encoding.Encode(route)
		if err != nil {
			return err
		}
	}
	return nil
}

var routesCommand = &cli.Command{
	Name:    "routes",
	Usage:   "Query routes for an athlete from Strava",
	Aliases: []string{"R"},
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "count",
			Aliases: []string{"N"},
			Value:   0,
			Usage:   "Count",
		},
	},
	Action: routes,
}

func entityWithArgs(c *cli.Context, f entityFunc, args []string) error {
	client, err := NewAPIClient(c)
	if err != nil {
		return err
	}
	for i := 0; i < len(args); i++ {
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		x, err := strconv.ParseInt(args[i], 0, 64)
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

func entity(c *cli.Context, f entityFunc) error {
	return entityWithArgs(c, f, c.Args().Slice())
}

var streamFlag = &cli.StringSliceFlag{
	Name:    "stream",
	Aliases: []string{"s"},
	Value:   cli.NewStringSlice(),
	Usage:   "Streams to include in the activity",
}

var activityCommand = &cli.Command{
	Name:    "activity",
	Aliases: []string{"a"},
	Usage:   "Query an activity from Strava",
	Flags:   []cli.Flag{streamFlag},
	Action: func(c *cli.Context) error {
		return entity(c, func(ctx context.Context, client *strava.Client, id int64) (interface{}, error) {
			return client.Activity.Activity(ctx, id, c.StringSlice("stream")...)
		})
	},
}

var streamsCommand = &cli.Command{
	Name:    "stream",
	Aliases: []string{"s"},
	Usage:   "Query streams for an activity from Strava",
	Flags:   []cli.Flag{streamFlag},
	Action: func(c *cli.Context) error {
		return entity(c, func(ctx context.Context, client *strava.Client, id int64) (interface{}, error) {
			streams := append([]string{"latlng", "altitude", "time"}, c.StringSlice("stream")...)
			return client.Activity.Streams(ctx, id, streams...)
		})
	},
}

var routeCommand = &cli.Command{
	Name:    "route",
	Aliases: []string{"r"},
	Usage:   "Query a route from Strava",
	Action: func(c *cli.Context) error {
		return entity(c, func(ctx context.Context, client *strava.Client, id int64) (interface{}, error) {
			return client.Route.Route(ctx, id)
		})
	},
}
