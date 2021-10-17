package cyclinganalytics

import (
	"context"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
	"golang.org/x/time/rate"

	api "github.com/bzimmer/activity"
	"github.com/bzimmer/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/activity"
)

const provider = "cyclinganalytics"

func athlete(c *cli.Context) error {
	client := pkg.Runtime(c).CyclingAnalytics
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	athlete, err := client.User.Me(ctx)
	if err != nil {
		return err
	}
	pkg.Runtime(c).Metrics.IncrCounter([]string{provider, c.Command.Name}, 1)
	return pkg.Runtime(c).Encoder.Encode(athlete)
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
	client := pkg.Runtime(c).CyclingAnalytics
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	rides, err := client.Rides.Rides(ctx, cyclinganalytics.Me, api.Pagination{Total: c.Int("count")})
	if err != nil {
		return err
	}
	enc := pkg.Runtime(c).Encoder
	for _, ride := range rides {
		if err := enc.Encode(ride); err != nil {
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
	client := pkg.Runtime(c).CyclingAnalytics
	opts := cyclinganalytics.RideOptions{
		Streams: []string{"latitude", "longitude", "elevation"},
	}
	enc := pkg.Runtime(c).Encoder
	args := c.Args()
	for i := 0; i < args.Len(); i++ {
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
		if err := enc.Encode(ride); err != nil {
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

func streamsetsCommand() *cli.Command {
	return &cli.Command{
		Name:  "streamsets",
		Usage: "Return the set of available streams for query",
		Action: func(c *cli.Context) error {
			client := pkg.Runtime(c).CyclingAnalytics
			if err := pkg.Runtime(c).Encoder.Encode(client.Rides.StreamSets()); err != nil {
				return err
			}
			return nil
		},
	}
}

func Before(c *cli.Context) error {
	client, err := cyclinganalytics.NewClient(
		cyclinganalytics.WithTokenCredentials(
			c.String("cyclinganalytics-access-token"), c.String("cyclinganalytics-refresh-token"), time.Time{}),
		cyclinganalytics.WithAutoRefresh(c.Context),
		cyclinganalytics.WithHTTPTracing(c.Bool("http-tracing")),
		cyclinganalytics.WithRateLimiter(rate.NewLimiter(
			rate.Every(c.Duration("rate-limit")), c.Int("rate-burst"))))
	if err != nil {
		return err
	}
	pkg.Runtime(c).CyclingAnalytics = client
	return nil
}

func Command() *cli.Command {
	return &cli.Command{
		Name:        "cyclinganalytics",
		Aliases:     []string{"ca"},
		Category:    "activity",
		Usage:       "Query CyclingAnalytics",
		Description: "Operations supported by the CyclingAnalytics API",
		Flags:       append(AuthFlags(), activity.RateLimitFlags...),
		Before:      Before,
		Subcommands: []*cli.Command{
			activitiesCommand(),
			athleteCommand(),
			oauthCommand(),
			rideCommand(),
			streamsetsCommand(),
			activity.UploadCommand(func(c *cli.Context) (api.Uploader, error) {
				return pkg.Runtime(c).CyclingAnalytics.Uploader(), nil
			}),
		},
	}
}

func AuthFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "cyclinganalytics-client-id",
			Usage:   "cyclinganalytics client id",
			EnvVars: []string{"CYCLINGANALYTICS_CLIENT_ID"},
		},
		&cli.StringFlag{
			Name:    "cyclinganalytics-client-secret",
			Usage:   "cyclinganalytics client secret",
			EnvVars: []string{"CYCLINGANALYTICS_CLIENT_SECRET"},
		},
		&cli.StringFlag{
			Name:    "cyclinganalytics-access-token",
			Usage:   "cyclinganalytics access token",
			EnvVars: []string{"CYCLINGANALYTICS_ACCESS_TOKEN"},
		},
	}
}
