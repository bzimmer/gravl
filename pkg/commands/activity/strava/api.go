package strava

import (
	"context"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
	"golang.org/x/time/rate"

	"github.com/bzimmer/gravl/pkg/commands"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/eval"
	"github.com/bzimmer/gravl/pkg/providers/activity"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

type entityFunc func(context.Context, *strava.Client, int64) (interface{}, error)

func NewAPIClient(c *cli.Context) (*strava.Client, error) {
	return strava.NewClient(
		strava.WithTokenCredentials(
			c.String("strava.access-token"), c.String("strava.refresh-token"), time.Now().Add(-1*time.Minute)),
		strava.WithClientCredentials(c.String("strava.client-id"), c.String("strava.client-secret")),
		strava.WithAutoRefresh(c.Context),
		strava.WithHTTPTracing(c.Bool("http-tracing")),
		strava.WithRateLimiter(rate.NewLimiter(
			rate.Every(c.Duration("rate-limit")), c.Int("rate-burst"))))
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

func filter(c *cli.Context) (func(ctx context.Context, act *strava.Activity) (bool, error), error) {
	f := func(ctx context.Context, act *strava.Activity) (bool, error) { return true, nil }
	if c.IsSet("filter") {
		var evaluator eval.Evaluator
		evaluator, err := commands.Evaluator(c.String("filter"))
		if err != nil {
			return nil, err
		}
		f = evaluator.Bool
	}
	return f, nil
}

func attributer(c *cli.Context) (func(ctx context.Context, act *strava.Activity) (interface{}, error), error) {
	f := func(ctx context.Context, act *strava.Activity) (interface{}, error) { return act, nil }
	if c.IsSet("attribute") {
		var evaluator eval.Evaluator
		evaluator, err := commands.Evaluator(c.String("attribute"))
		if err != nil {
			return nil, err
		}
		f = evaluator.Eval
	}
	return f, nil
}

func activities(c *cli.Context) error {
	client, err := NewAPIClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()

	f, err := filter(c)
	if err != nil {
		return err
	}
	g, err := attributer(c)
	if err != nil {
		return err
	}

	var ok bool
	var res *strava.ActivityResult
	acts := client.Activity.Activities(ctx, activity.Pagination{Total: c.Int("count")})
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case res, ok = <-acts:
			if !ok {
				return nil
			}
			if res.Err != nil {
				return res.Err
			}
			// filter
			ok, err = f(ctx, res.Activity)
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
			// extract
			res, err := g(ctx, res.Activity)
			if err != nil {
				return err
			}
			// encode
			if err = encoding.Encode(res); err != nil {
				return err
			}
		}
	}
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
			Usage:   "The number of activities to query from Strava (the number returned will be <= N)",
		},
		&cli.StringFlag{
			Name:    "filter",
			Aliases: []string{"f"},
			Usage:   "Expression for filtering activities to remove",
		},
		&cli.StringSliceFlag{
			Name:    "attribute",
			Aliases: []string{"B"},
			Usage:   "Evaluate the expression on an activity and return only those results",
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
			Usage:   "The number of routes to query from Strava (the number returned will be <= N)",
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
	Name:      "activity",
	Aliases:   []string{"a"},
	Usage:     "Query an activity from Strava",
	ArgsUsage: "ACTIVITY_ID (...)",
	Flags:     []cli.Flag{streamFlag},
	Action: func(c *cli.Context) error {
		return entity(c, func(ctx context.Context, client *strava.Client, id int64) (interface{}, error) {
			return client.Activity.Activity(ctx, id, c.StringSlice("stream")...)
		})
	},
}

var streamsetsCommand = &cli.Command{
	Name:  "streamsets",
	Usage: "Return the set of available streams for query",
	Action: func(c *cli.Context) error {
		client, err := NewAPIClient(c)
		if err != nil {
			return err
		}
		if err := encoding.Encode(client.Activity.StreamSets()); err != nil {
			return err
		}
		return nil
	},
}

var streamsCommand = &cli.Command{
	Name:      "streams",
	Aliases:   []string{"s"},
	Usage:     "Query streams for an activity from Strava",
	ArgsUsage: "ACTIVITY_ID (...)",
	Flags:     []cli.Flag{streamFlag},
	Action: func(c *cli.Context) error {
		return entity(c, func(ctx context.Context, client *strava.Client, id int64) (interface{}, error) {
			streams := append([]string{"latlng", "altitude", "time"}, c.StringSlice("stream")...)
			return client.Activity.Streams(ctx, id, streams...)
		})
	},
}

var routeCommand = &cli.Command{
	Name:      "route",
	Aliases:   []string{"r"},
	Usage:     "Query a route from Strava",
	ArgsUsage: "ROUTE_ID (...)",
	Action: func(c *cli.Context) error {
		return entity(c, func(ctx context.Context, client *strava.Client, id int64) (interface{}, error) {
			return client.Route.Route(ctx, id)
		})
	},
}

var photosCommand = &cli.Command{
	Name:      "photos",
	Aliases:   []string{""},
	Usage:     "Query photos from Strava",
	ArgsUsage: "ACTIVITY_ID (...)",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "size",
			Aliases: []string{"s"},
			Value:   2048,
		},
	},
	Action: func(c *cli.Context) error {
		return entity(c, func(ctx context.Context, client *strava.Client, id int64) (interface{}, error) {
			return client.Activity.Photos(ctx, id, c.Int("size"))
		})
	},
}
