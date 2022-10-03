package strava

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	api "github.com/bzimmer/activity"
	"github.com/bzimmer/activity/strava"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/activity"
	"github.com/bzimmer/gravl/eval"
)

const Provider = "strava"

var before sync.Once //nolint:gochecknoglobals

type entityFunc func(context.Context, *strava.Client, int64) (any, error)

func athlete(c *cli.Context) error {
	client := gravl.Runtime(c).Strava
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	athlete, err := client.Athlete.Athlete(ctx)
	if err != nil {
		return err
	}
	log.Info().Int64("id", int64(athlete.ID)).Str("username", athlete.Username).Msg(c.Command.Name)
	gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, c.Command.Name}, 1)
	return gravl.Runtime(c).Encoder.Encode(athlete)
}

func athleteCommand() *cli.Command {
	return &cli.Command{
		Name:    "athlete",
		Usage:   "Query an athlete from Strava",
		Aliases: []string{"t"},
		Action:  athlete,
	}
}

func refresh(c *cli.Context) error {
	client := gravl.Runtime(c).Strava
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	tokens, err := client.Auth.Refresh(ctx)
	if err != nil {
		return err
	}
	return gravl.Runtime(c).Encoder.Encode(tokens)
}

func refreshCommand() *cli.Command {
	return &cli.Command{
		Name:   "refresh",
		Usage:  "Acquire a new refresh token",
		Action: refresh,
	}
}

func evaluator(c *cli.Context, evaluation string) (eval.Evaluator, error) {
	if c.IsSet(evaluation) {
		var ev eval.Evaluator
		ev, err := gravl.Runtime(c).Evaluator(c.String(evaluation))
		if err != nil {
			return nil, err
		}
		return ev, nil
	}
	return nil, nil
}

func filter(c *cli.Context) (func(ctx context.Context, act *strava.Activity) (bool, error), error) {
	ev, err := evaluator(c, "filter")
	if err != nil {
		return nil, err
	}
	if ev == nil {
		return func(ctx context.Context, act *strava.Activity) (bool, error) { return true, nil }, nil
	}
	return ev.Bool, nil
}

func attributer(c *cli.Context) (func(ctx context.Context, act *strava.Activity) (any, error), error) {
	ev, err := evaluator(c, "attribute")
	if err != nil {
		return nil, err
	}
	if ev == nil {
		return func(ctx context.Context, act *strava.Activity) (any, error) { return act, nil }, nil
	}
	return ev.Eval, nil
}

func daterange(c *cli.Context) (strava.APIOption, error) {
	before, after, err := activity.DateRange(c)
	if err != nil {
		return nil, err
	}
	log.Info().Time("before", before).Time("after", after).Msg("date range")
	return strava.WithDateRange(before, after), nil
}

func activities(c *cli.Context) error {
	client := gravl.Runtime(c).Strava
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

	opt, err := daterange(c)
	if err != nil {
		return err
	}

	enc := gravl.Runtime(c).Encoder
	met := gravl.Runtime(c).Metrics

	met.IncrCounter([]string{Provider, c.Command.Name}, 1)
	defer func(t time.Time) {
		met.AddSample([]string{Provider, c.Command.Name}, float32(time.Since(t).Seconds()))
	}(time.Now())

	acts := client.Activity.Activities(ctx, api.Pagination{Total: c.Int("count")}, opt)
	return strava.ActivitiesIter(acts, func(act *strava.Activity) (bool, error) {
		// filter
		var ok bool
		ok, err = f(ctx, act)
		if err != nil {
			return false, err
		}
		if !ok {
			return true, nil
		}
		// extract
		var ext any
		ext, err = g(ctx, act)
		if err != nil {
			return false, err
		}
		met.IncrCounter([]string{Provider, "activity"}, 1)
		log.Info().
			Time("date", act.StartDateLocal).
			Int64("id", act.ID).
			Str("name", act.Name).
			Str("type", act.Type).
			Msg(c.Command.Name)
		// encode
		if err = enc.Encode(ext); err != nil {
			return false, err
		}
		return true, nil
	})
}

func activitiesCommand() *cli.Command {
	return &cli.Command{
		Name:    "activities",
		Usage:   "Query activities for an athlete from Strava",
		Aliases: []string{"A"},
		Flags: append([]cli.Flag{
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
		}, activity.DateRangeFlags()...),
		Action: activities,
	}
}

func routes(c *cli.Context) error {
	client := gravl.Runtime(c).Strava
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	athlete, err := client.Athlete.Athlete(ctx)
	if err != nil {
		return err
	}
	routes, err := client.Route.Routes(ctx, athlete.ID, api.Pagination{Total: c.Int("count")})
	if err != nil {
		return err
	}
	enc := gravl.Runtime(c).Encoder
	met := gravl.Runtime(c).Metrics
	met.IncrCounter([]string{Provider, c.Command.Name}, 1)
	for _, route := range routes {
		met.IncrCounter([]string{Provider, "route"}, 1)
		if err = enc.Encode(route); err != nil {
			return err
		}
	}
	return nil
}

func routesCommand() *cli.Command {
	return &cli.Command{
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
}

func entityWithArgs(c *cli.Context, f entityFunc, args []string) error { //nolint:gocognit
	if len(args) == 0 {
		log.Info().Str("entity", c.Command.Name).Msg("no arguments provided")
		return nil
	}
	enc := gravl.Runtime(c).Encoder
	met := gravl.Runtime(c).Metrics
	client := gravl.Runtime(c).Strava

	concurrency := c.Int("concurrency")
	if len(args) < concurrency {
		concurrency = len(args)
	}
	if concurrency <= 0 {
		concurrency = 1
	}

	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()

	argc := make(chan int64)
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		defer close(argc)
		for _, arg := range args {
			x, err := strconv.ParseInt(arg, 0, 64)
			if err != nil {
				return err
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case argc <- x:
			}
		}
		return nil
	})
	for i := 0; i < concurrency; i++ {
		grp.Go(func() error {
			for x := range argc {
				t := time.Now()
				log.Info().Int64("id", x).Str("command", c.Command.Name).Msg("executing")
				v, err := f(ctx, client, x)
				if err != nil {
					return err
				}
				met.IncrCounter([]string{Provider, c.Command.Name}, 1)
				if err = enc.Encode(v); err != nil {
					return err
				}
				met.AddSample([]string{Provider, c.Command.Name}, float32(time.Since(t).Seconds()))
			}
			return nil
		})
	}
	return grp.Wait()
}

func entity(c *cli.Context, f entityFunc) error {
	return entityWithArgs(c, f, c.Args().Slice())
}

func streamFlag(streams ...string) cli.Flag {
	return &cli.StringSliceFlag{
		Name:    "stream",
		Aliases: []string{"s"},
		Value:   cli.NewStringSlice(streams...),
		Usage:   "Streams to include in the activity",
	}
}

func activityCommand() *cli.Command {
	return &cli.Command{
		Name:      "activity",
		Aliases:   []string{"a"},
		Usage:     "Query an activity from Strava",
		ArgsUsage: "ACTIVITY_ID (...)",
		Flags:     []cli.Flag{streamFlag()},
		Action: func(c *cli.Context) error {
			s := make(map[string]bool)
			for _, x := range c.StringSlice("stream") {
				s[x] = true
			}
			var streams []string
			for stream := range s {
				streams = append(streams, stream)
			}
			return entity(c, func(ctx context.Context, client *strava.Client, id int64) (any, error) {
				act, err := client.Activity.Activity(ctx, id, streams...)
				if err != nil {
					return nil, err
				}
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
}

func updateFlags() []cli.Flag {
	var flags []cli.Flag
	var pairs = [][]string{
		{"name", "Set the name for the activity"},
		{"gear", "Set the gear id for the activity"},
		{"sport", "Set the sport for the activity"},
		{"description", "Set the description for the activity"},
	}
	for i := range pairs {
		name := pairs[i][0]
		usage := pairs[i][1]
		flags = append(flags, &cli.StringFlag{Name: name, Usage: usage})
	}
	pairs = [][]string{
		{"hidden", "Hide the activity from the home dashboard"},
		{"no-hidden", "Display the activity on the home dashboard"},
		{"commute", "The activity is a commute"},
		{"no-commute", "The activity is not a commute"},
		{"trainer", "The activity was completed on a trainer"},
		{"no-trainer", "The activity was not completed on a trainer"},
	}
	for i := range pairs {
		name := pairs[i][0]
		usage := pairs[i][1]
		flags = append(flags, &cli.BoolFlag{Name: name, Usage: usage})
	}
	return flags
}

func updateCommand() *cli.Command { //nolint:gocognit
	return &cli.Command{
		Name:      "update",
		ArgsUsage: "ACTIVITY_ID (...)",
		Flags:     updateFlags(),
		Action: func(c *cli.Context) error {
			met := gravl.Runtime(c).Metrics
			return entity(c, func(ctx context.Context, client *strava.Client, id int64) (any, error) {
				update := &strava.UpdatableActivity{ID: id}
				if c.IsSet("name") {
					val := c.String("name")
					update.Name = &val
					met.IncrCounter([]string{Provider, c.Command.Name, "name"}, 1)
				}
				if c.IsSet("sport") {
					val := c.String("sport")
					update.SportType = &val
					met.IncrCounter([]string{Provider, c.Command.Name, "sport"}, 1)
				}
				if c.IsSet("gear") {
					val := c.String("gear")
					update.GearID = &val
					met.IncrCounter([]string{Provider, c.Command.Name, "gear"}, 1)
				}
				if c.IsSet("description") {
					val := c.String("description")
					update.Description = &val
					met.IncrCounter([]string{Provider, c.Command.Name, "description"}, 1)
				}
				switch {
				case c.IsSet("hidden") && c.IsSet("no-hidden"):
					return nil, errors.New("only one of hidden or no-hidden can be specified")
				case c.IsSet("hidden"):
					val := true
					update.Hidden = &val
					met.IncrCounter([]string{Provider, c.Command.Name, "hidden"}, 1)
				case c.IsSet("no-hidden"):
					val := false
					update.Hidden = &val
					met.IncrCounter([]string{Provider, c.Command.Name, "no-hidden"}, 1)
				}
				switch {
				case c.IsSet("commute") && c.IsSet("no-commute"):
					return nil, errors.New("only one of commute or no-commute can be specified")
				case c.IsSet("commute"):
					val := true
					update.Commute = &val
					met.IncrCounter([]string{Provider, c.Command.Name, "commute"}, 1)
				case c.IsSet("no-commute"):
					val := false
					update.Commute = &val
					met.IncrCounter([]string{Provider, c.Command.Name, "no-commute"}, 1)
				}
				switch {
				case c.IsSet("trainer") && c.IsSet("no-trainer"):
					return nil, errors.New("only one of trainer or no-trainer can be specified")
				case c.IsSet("trainer"):
					val := true
					update.Trainer = &val
					met.IncrCounter([]string{Provider, c.Command.Name, "trainer"}, 1)
				case c.IsSet("no-trainer"):
					val := false
					update.Trainer = &val
					met.IncrCounter([]string{Provider, c.Command.Name, "no-trainer"}, 1)
				}
				return client.Activity.Update(ctx, update)
			})
		},
	}
}

func streamsCommand() *cli.Command {
	return &cli.Command{
		Name:      "streams",
		Aliases:   []string{"s"},
		Usage:     "Query streams for an activity from Strava",
		ArgsUsage: "ACTIVITY_ID (...)",
		Flags:     []cli.Flag{streamFlag("latlng", "altitude", "time")},
		Action: func(c *cli.Context) error {
			s := make(map[string]bool)
			for _, x := range c.StringSlice("stream") {
				s[x] = true
			}
			var streams []string
			for stream := range s {
				streams = append(streams, stream)
			}
			log.Info().Strs("streams", streams).Msg(c.Command.Name)
			return entity(c, func(ctx context.Context, client *strava.Client, id int64) (any, error) {
				return client.Activity.Streams(ctx, id, streams...)
			})
		},
	}
}

func routeCommand() *cli.Command {
	return &cli.Command{
		Name:      "route",
		Aliases:   []string{"r"},
		Usage:     "Query a route from Strava",
		ArgsUsage: "ROUTE_ID (...)",
		Action: func(c *cli.Context) error {
			return entity(c, func(ctx context.Context, client *strava.Client, id int64) (any, error) {
				return client.Route.Route(ctx, id)
			})
		},
	}
}

func photosCommand() *cli.Command {
	return &cli.Command{
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
			return entity(c, func(ctx context.Context, client *strava.Client, id int64) (any, error) {
				return client.Activity.Photos(ctx, id, c.Int("size"))
			})
		},
	}
}

func streamSetsCommand() *cli.Command {
	return &cli.Command{
		Name:  "streamsets",
		Usage: "Return the set of available streams for query",
		Action: func(c *cli.Context) error {
			client := gravl.Runtime(c).Strava
			gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, c.Command.Name}, 1)
			if err := gravl.Runtime(c).Encoder.Encode(client.Activity.StreamSets()); err != nil {
				return err
			}
			return nil
		},
	}
}

func oauthCommand() *cli.Command {
	return activity.OAuthCommand(&activity.OAuthConfig{
		Port:     9001,
		Provider: Provider,
		Scopes:   []string{"read_all,profile:read_all,activity:read_all,activity:write"},
	})
}

func Before(c *cli.Context) error {
	var err error
	before.Do(func() {
		var client *strava.Client
		client, err = strava.NewClient(
			strava.WithTokenCredentials(
				c.String("strava-refresh-token"), c.String("strava-refresh-token"), time.Now().Add(-1*time.Minute)),
			strava.WithClientCredentials(c.String("strava-client-id"), c.String("strava-client-secret")),
			strava.WithAutoRefresh(c.Context),
			strava.WithHTTPTracing(c.Bool("http-tracing")),
			strava.WithRateLimiter(rate.NewLimiter(
				rate.Every(c.Duration("rate-limit")), c.Int("rate-burst"))))
		if err != nil {
			return
		}
		gravl.Runtime(c).Endpoints[Provider] = strava.Endpoint()
		gravl.Runtime(c).Strava = client
		gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, "client", "created"}, 1)
		log.Info().Msg("created strava client")
	})
	return err
}

func Command() *cli.Command {
	return &cli.Command{
		Name:        Provider,
		Category:    "activity",
		Usage:       "Query Strava for rides and routes",
		Description: "Operations supported by the Strava API",
		Flags:       append(AuthFlags(), activity.RateLimitFlags()...),
		Before:      Before,
		Subcommands: []*cli.Command{
			activitiesCommand(),
			activityCommand(),
			athleteCommand(),
			oauthCommand(),
			photosCommand(),
			refreshCommand(),
			routeCommand(),
			routesCommand(),
			streamsCommand(),
			streamSetsCommand(),
			updateCommand(),
			webhookCommand(),
		},
	}
}

func AuthFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "strava-client-id",
			Usage:   "Strava client id",
			EnvVars: []string{"STRAVA_CLIENT_ID"},
		},
		&cli.StringFlag{
			Name:    "strava-client-secret",
			Usage:   "Strava client secret",
			EnvVars: []string{"STRAVA_CLIENT_SECRET"},
		},
		&cli.StringFlag{
			Name:    "strava-refresh-token",
			Usage:   "Strava refresh token",
			EnvVars: []string{"STRAVA_REFRESH_TOKEN"},
		},
	}
}
