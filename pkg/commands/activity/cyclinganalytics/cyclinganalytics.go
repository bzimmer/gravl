package cyclinganalytics

import (
	"context"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"golang.org/x/time/rate"

	actcmd "github.com/bzimmer/gravl/pkg/commands/activity"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/activity"
	"github.com/bzimmer/gravl/pkg/providers/activity/cyclinganalytics"
)

func NewClient(c *cli.Context) (*cyclinganalytics.Client, error) {
	return cyclinganalytics.NewClient(
		cyclinganalytics.WithTokenCredentials(
			c.String("cyclinganalytics.access-token"), c.String("cyclinganalytics.refresh-token"), time.Time{}),
		cyclinganalytics.WithAutoRefresh(c.Context),
		cyclinganalytics.WithHTTPTracing(c.Bool("http-tracing")),
		cyclinganalytics.WithRateLimiter(rate.NewLimiter(
			rate.Every(c.Duration("rate-limit")), c.Int("rate-burst"))))
}

func athlete(c *cli.Context) error {
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	athlete, err := client.User.Me(ctx)
	if err != nil {
		return err
	}
	return encoding.For(c).Encode(athlete)
}

var athleteCommand = &cli.Command{
	Name:    "athlete",
	Aliases: []string{"t"},
	Usage:   "Query for the authenticated athlete",
	Action:  athlete,
}

func activities(c *cli.Context) error {
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	rides, err := client.Rides.Rides(ctx, cyclinganalytics.Me, activity.Pagination{Total: c.Int("count")})
	if err != nil {
		return err
	}
	enc := encoding.For(c)
	for _, ride := range rides {
		if err := enc.Encode(ride); err != nil {
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
			Usage:   "The number of activities to query from CA (the number returned will be <= N)",
		},
	},
	Action: activities,
}

func ride(c *cli.Context) error {
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	opts := cyclinganalytics.RideOptions{
		Streams: []string{"latitude", "longitude", "elevation"},
	}
	enc := encoding.For(c)
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
		if err = enc.Encode(ride); err != nil {
			return err
		}
	}
	return nil
}

var rideCommand = &cli.Command{
	Name:    "activity",
	Aliases: []string{"a"},
	Usage:   "Query an activity for the authenticated athlete",
	Action:  ride,
}

var streamsetsCommand = &cli.Command{
	Name:  "streamsets",
	Usage: "Return the set of available streams for query",
	Action: func(c *cli.Context) error {
		client, err := NewClient(c)
		if err != nil {
			return err
		}
		if err := encoding.For(c).Encode(client.Rides.StreamSets()); err != nil {
			return err
		}
		return nil
	},
}

var Command = &cli.Command{
	Name:        "cyclinganalytics",
	Aliases:     []string{"ca"},
	Category:    "activity",
	Usage:       "Query CyclingAnalytics",
	Description: "Operations supported by the CyclingAnalytics API",
	Flags:       append(AuthFlags, actcmd.RateLimitFlags...),
	Subcommands: []*cli.Command{
		activitiesCommand,
		athleteCommand,
		oauthCommand,
		rideCommand,
		streamsetsCommand,
		actcmd.UploadCommand(func(c *cli.Context) (activity.Uploader, error) {
			client, err := NewClient(c)
			if err != nil {
				return nil, err
			}
			return client.Uploader(), nil
		}),
	},
}

var AuthFlags = []cli.Flag{
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "cyclinganalytics.client-id",
		Usage: "API key for CyclingAnalytics API",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "cyclinganalytics.client-secret",
		Usage: "API secret for CyclingAnalytics API",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "cyclinganalytics.access-token",
		Usage: "Access token for CyclingAnalytics API",
	}),
}
