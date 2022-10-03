package cyclinganalytics

import (
	"context"
	"strconv"
	"sync"
	"time"

	api "github.com/bzimmer/activity"
	"github.com/bzimmer/activity/cyclinganalytics"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/time/rate"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/activity"
)

const Provider = "cyclinganalytics"

var before sync.Once //nolint:gochecknoglobals

func athlete(c *cli.Context) error {
	client := gravl.Runtime(c).CyclingAnalytics
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	athlete, err := client.User.Me(ctx)
	if err != nil {
		return err
	}
	log.Info().Int64("id", int64(athlete.ID)).Str("username", athlete.Email).Msg(c.Command.Name)
	gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, c.Command.Name}, 1)
	return gravl.Runtime(c).Encoder.Encode(athlete)
}

func athleteCommand() *cli.Command {
	return &cli.Command{
		Name:    "athlete",
		Aliases: []string{"t"},
		Usage:   "Query for the authenticated athlete",
		Action:  athlete,
	}
}

func activities(c *cli.Context) error {
	client := gravl.Runtime(c).CyclingAnalytics
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	rides, err := client.Rides.Rides(ctx, cyclinganalytics.Me, api.Pagination{Total: c.Int("count")})
	if err != nil {
		return err
	}
	gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, c.Command.Name}, 1)
	enc := gravl.Runtime(c).Encoder
	for _, ride := range rides {
		gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, "activity"}, 1)
		if err = enc.Encode(ride); err != nil {
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
				Usage:   "The number of activities to query from CA (the number returned will be <= N)",
			},
		},
		Action: activities,
	}
}

func ride(c *cli.Context) error {
	client := gravl.Runtime(c).CyclingAnalytics
	opts := cyclinganalytics.WithRideOptions(cyclinganalytics.RideOptions{
		Streams: []string{"latitude", "longitude", "elevation"},
	})
	enc := gravl.Runtime(c).Encoder
	args := c.Args()
	for i := 0; i < args.Len(); i++ {
		err := func() error {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			rideID, err := strconv.ParseInt(args.Get(i), 0, 64)
			if err != nil {
				return err
			}
			ride, err := client.Rides.Ride(ctx, rideID, opts)
			if err != nil {
				return err
			}
			gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, c.Command.Name}, 1)
			if err = enc.Encode(ride); err != nil {
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

func rideCommand() *cli.Command {
	return &cli.Command{
		Name:    "activity",
		Aliases: []string{"a"},
		Usage:   "Query an activity for the authenticated athlete",
		Action:  ride,
	}
}

func streamSetsCommand() *cli.Command {
	return &cli.Command{
		Name:  "streamsets",
		Usage: "Return the set of available streams for query",
		Action: func(c *cli.Context) error {
			client := gravl.Runtime(c).CyclingAnalytics
			gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, c.Command.Name}, 1)
			if err := gravl.Runtime(c).Encoder.Encode(client.Rides.StreamSets()); err != nil {
				return err
			}
			return nil
		},
	}
}

func oauthCommand() *cli.Command {
	return activity.OAuthCommand(&activity.OAuthConfig{
		Port:     9002,
		Provider: Provider,
		Scopes:   []string{"read_account,read_email,read_athlete,read_rides,create_rides"},
	})
}

func Before(c *cli.Context) error {
	var err error
	before.Do(func() {
		var client *cyclinganalytics.Client
		client, err = cyclinganalytics.NewClient(
			cyclinganalytics.WithTokenCredentials(
				c.String("cyclinganalytics-access-token"), c.String("cyclinganalytics-refresh-token"), time.Time{}),
			cyclinganalytics.WithAutoRefresh(c.Context),
			cyclinganalytics.WithHTTPTracing(c.Bool("http-tracing")),
			cyclinganalytics.WithRateLimiter(rate.NewLimiter(
				rate.Every(c.Duration("rate-limit")), c.Int("rate-burst"))))
		if err != nil {
			return
		}
		gravl.Runtime(c).Endpoints[Provider] = cyclinganalytics.Endpoint()
		gravl.Runtime(c).CyclingAnalytics = client
		gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, "client", "created"}, 1)
		log.Info().Msg("created cyclinganalytics client")
	})
	return err
}

func Command() *cli.Command {
	return &cli.Command{
		Name:        "cyclinganalytics",
		Aliases:     []string{"ca"},
		Category:    "activity",
		Usage:       "Query CyclingAnalytics",
		Description: "Operations supported by the CyclingAnalytics API",
		Flags:       append(AuthFlags(), activity.RateLimitFlags()...),
		Before:      Before,
		Subcommands: []*cli.Command{
			activitiesCommand(),
			athleteCommand(),
			oauthCommand(),
			rideCommand(),
			streamSetsCommand(),
		},
	}
}

func AuthFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "cyclinganalytics-client-id",
			Usage:   "CyclingAnalytics client id",
			EnvVars: []string{"CYCLINGANALYTICS_CLIENT_ID"},
		},
		&cli.StringFlag{
			Name:    "cyclinganalytics-client-secret",
			Usage:   "CyclingAnalytics client secret",
			EnvVars: []string{"CYCLINGANALYTICS_CLIENT_SECRET"},
		},
		&cli.StringFlag{
			Name:    "cyclinganalytics-access-token",
			Usage:   "CyclingAnalytics access token",
			EnvVars: []string{"CYCLINGANALYTICS_ACCESS_TOKEN"},
		},
	}
}
