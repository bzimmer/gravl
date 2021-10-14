package strava

import (
	"context"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"

	"github.com/bzimmer/activity"
	"github.com/bzimmer/activity/strava"
	"github.com/bzimmer/gravl/pkg"
	actcmd "github.com/bzimmer/gravl/pkg/activity"
	"github.com/bzimmer/gravl/pkg/eval"
)

const provider = "strava"

type entityFunc func(context.Context, *strava.Client, int64) (interface{}, error)

func athlete(c *cli.Context) error {
	client := pkg.Runtime(c).Strava
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	athlete, err := client.Athlete.Athlete(ctx)
	if err != nil {
		return err
	}
	pkg.Runtime(c).Metrics.IncrCounter([]string{provider, c.Command.Name}, 1)
	return pkg.Runtime(c).Encoder.Encode(athlete)
}

var athleteCommand = &cli.Command{
	Name:    "athlete",
	Usage:   "Query an athlete from Strava",
	Aliases: []string{"t"},
	Action:  athlete,
}

func refresh(c *cli.Context) error {
	client := pkg.Runtime(c).Strava
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	tokens, err := client.Auth.Refresh(ctx)
	if err != nil {
		return err
	}
	return pkg.Runtime(c).Encoder.Encode(tokens)
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
		evaluator, err := pkg.Runtime(c).Evaluator(c.String("filter"))
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
		evaluator, err := pkg.Runtime(c).Evaluator(c.String("attribute"))
		if err != nil {
			return nil, err
		}
		f = evaluator.Eval
	}
	return f, nil
}

func activities(c *cli.Context) error {
	client := pkg.Runtime(c).Strava
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

	enc := pkg.Runtime(c).Encoder
	met := pkg.Runtime(c).Metrics

	met.IncrCounter([]string{provider, c.Command.Name}, 1)
	defer func(t time.Time) {
		met.AddSample([]string{provider, c.Command.Name}, float32(time.Since(t).Seconds()))
	}(time.Now())

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
			act, err := g(ctx, res.Activity)
			if err != nil {
				return err
			}
			met.IncrCounter([]string{provider, "activity"}, 1)
			log.Info().
				Time("date", res.Activity.StartDateLocal).
				Int64("id", res.Activity.ID).
				Str("name", res.Activity.Name).
				Str("type", res.Activity.Type).
				Msg(c.Command.Name)
			// encode
			if err := enc.Encode(act); err != nil {
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
	client := pkg.Runtime(c).Strava
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
	enc := pkg.Runtime(c).Encoder
	for _, route := range routes {
		err = enc.Encode(route)
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
	client := pkg.Runtime(c).Strava
	argc := make(chan string, c.NArg())
	go func() {
		defer close(argc)
		for _, arg := range args {
			argc <- arg
		}
	}()

	concurrency := 5
	if len(args) < concurrency {
		concurrency = len(args)
	}

	enc := pkg.Runtime(c).Encoder
	met := pkg.Runtime(c).Metrics
	defer func(t time.Time) {
		met.AddSample([]string{provider, c.Command.Name}, float32(time.Since(t).Seconds()))
	}(time.Now())

	grp, ctx := errgroup.WithContext(c.Context)
	for i := 0; i < concurrency; i++ {
		grp.Go(func() error {
			gtx, cancel := context.WithTimeout(ctx, c.Duration("timeout"))
			defer cancel()
			for arg := range argc {
				x, err := strconv.ParseInt(arg, 0, 64)
				if err != nil {
					return err
				}
				v, err := f(gtx, client, x)
				if err != nil {
					return err
				}
				if err := enc.Encode(v); err != nil {
					return err
				}
			}
			return nil
		})
	}
	return grp.Wait()
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
			act, err := client.Activity.Activity(ctx, id, c.StringSlice("stream")...)
			if err != nil {
				return nil, err
			}
			pkg.Runtime(c).Metrics.IncrCounter([]string{provider, c.Command.Name}, 1)
			log.Info().
				Time("date", act.StartDateLocal).
				Int64("id", act.ID).
				Str("name", act.Name).
				Str("type", act.Type).
				Msg("activity")
			return act, nil
		})
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
			log.Info().Int64("id", id).Msg("querying streams")
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
			log.Info().Int64("id", id).Msg("querying route")
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
			defer func(t time.Time) {
				log.Info().Int64("id", id).Dur("elapsed", time.Since(t)).Msg("querying photos")
			}(time.Now())
			return client.Activity.Photos(ctx, id, c.Int("size"))
		})
	},
}

var streamsetsCommand = &cli.Command{
	Name:  "streamsets",
	Usage: "Return the set of available streams for query",
	Action: func(c *cli.Context) error {
		client := pkg.Runtime(c).Strava
		pkg.Runtime(c).Metrics.IncrCounter([]string{provider, c.Command.Name}, 1)
		if err := pkg.Runtime(c).Encoder.Encode(client.Activity.StreamSets()); err != nil {
			return err
		}
		return nil
	},
}

func Before(c *cli.Context) error {
	client, err := strava.NewClient(
		strava.WithTokenCredentials(
			c.String("strava-refresh-token"), c.String("strava-refresh-token"), time.Now().Add(-1*time.Minute)),
		strava.WithClientCredentials(c.String("strava-client-id"), c.String("strava-client-secret")),
		strava.WithAutoRefresh(c.Context),
		strava.WithHTTPTracing(c.Bool("http-tracing")),
		strava.WithRateLimiter(rate.NewLimiter(
			rate.Every(c.Duration("rate-limit")), c.Int("rate-burst"))))
	if err != nil {
		return err
	}
	pkg.Runtime(c).Strava = client
	return nil
}

func Command() *cli.Command {
	return &cli.Command{
		Name:        provider,
		Category:    "activity",
		Usage:       "Query Strava for rides and routes",
		Description: "Operations supported by the Strava API",
		Flags:       append(AuthFlags, actcmd.RateLimitFlags...),
		Before:      Before,
		Subcommands: []*cli.Command{
			activitiesCommand,
			activityCommand,
			athleteCommand,
			oauthCommand,
			photosCommand,
			refreshCommand,
			routeCommand,
			routesCommand,
			streamsCommand,
			streamsetsCommand,
			actcmd.UploadCommand(func(c *cli.Context) (activity.Uploader, error) {
				return pkg.Runtime(c).Strava.Uploader(), nil
			}),
			webhookCommand,
		},
	}
}

var AuthFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "strava-client-id",
		Required: true,
		Usage:    "strava client id",
		EnvVars:  []string{"STRAVA_CLIENT_ID"},
	},
	&cli.StringFlag{
		Name:     "strava-client-secret",
		Required: true,
		Usage:    "strava client secret",
		EnvVars:  []string{"STRAVA_CLIENT_SECRET"},
	},
	&cli.StringFlag{
		Name:     "strava-refresh-token",
		Required: true,
		Usage:    "strava refresh token",
		EnvVars:  []string{"STRAVA_REFRESH_TOKEN"},
	},
}
