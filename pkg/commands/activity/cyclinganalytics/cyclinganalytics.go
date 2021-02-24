package cyclinganalytics

import (
	"context"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/commands/activity/internal"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/activity"
	"github.com/bzimmer/gravl/pkg/providers/activity/cyclinganalytics"
)

func NewClient(c *cli.Context) (*cyclinganalytics.Client, error) {
	return cyclinganalytics.NewClient(
		cyclinganalytics.WithTokenCredentials(
			c.String("cyclinganalytics.access-token"), c.String("cyclinganalytics.refresh-token"), time.Time{}),
		cyclinganalytics.WithAutoRefresh(c.Context),
		cyclinganalytics.WithHTTPTracing(c.Bool("http-tracing")))
}

func poll(ctx context.Context, client *cyclinganalytics.Client, uploadID int64, follow bool) error {
	pc := client.Rides.Poll(ctx, uploadID)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case res, ok := <-pc:
			if !ok {
				return nil
			}
			if res.Err != nil {
				return res.Err
			}
			if err := encoding.Encode(res.Upload); err != nil {
				return err
			}
		}
		if !follow {
			return nil
		}
	}
}

func upload(c *cli.Context) error {
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	args := c.Args()
	dryrun := c.Bool("dryrun")
	for i := 0; i < args.Len(); i++ {
		files, err := internal.Collect(args.Get(i), nil)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			log.Warn().Msg("no files specified")
			return nil
		}
		for _, file := range files {
			defer file.Close()
			if dryrun {
				log.Info().Str("file", file.Name).Bool("dryrun", dryrun).Msg("uploading")
				continue
			}
			log.Info().Str("file", file.Name).Msg("uploading")
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			u, err := client.Rides.Upload(ctx, file)
			if err != nil {
				return err
			}
			if !c.Bool("poll") {
				return encoding.Encode(u)
			}
			if err := poll(ctx, client, u.ID, true); err != nil {
				return err
			}
		}
	}
	return nil
}

func status(c *cli.Context) error {
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	args := c.Args()
	for i := 0; i < args.Len(); i++ {
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		uploadID, err := strconv.ParseInt(args.Get(i), 0, 64)
		if err != nil {
			return err
		}
		if err := poll(ctx, client, uploadID, c.Bool("poll")); err != nil {
			return err
		}
	}
	return nil
}

var uploadCommand = &cli.Command{
	Name:      "upload",
	Aliases:   []string{"u"},
	Usage:     "Upload an activity file",
	ArgsUsage: "{FILE | DIRECTORY}",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "status",
			Aliases: []string{"s"},
			Value:   false,
			Usage:   "Check the status of the upload",
		},
		&cli.BoolFlag{
			Name:    "poll",
			Aliases: []string{"p"},
			Value:   false,
			Usage:   "Continually check the status of the request until it is completed",
		},
		&cli.BoolFlag{
			Name:    "dryrun",
			Aliases: []string{"n"},
			Value:   false,
			Usage:   "Show the files which would be uploaded but do not upload them",
		},
	},
	Action: func(c *cli.Context) error {
		if c.Bool("status") {
			return status(c)
		}
		return upload(c)
	},
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
	return encoding.Encode(athlete)
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
	for _, ride := range rides {
		if err := encoding.Encode(ride); err != nil {
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
		if err = encoding.Encode(ride); err != nil {
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

var Command = &cli.Command{
	Name:        "cyclinganalytics",
	Aliases:     []string{"ca"},
	Category:    "activity",
	Usage:       "Query CyclingAnalytics",
	Description: "Operations supported by the Cycling Analytics website",
	Flags:       AuthFlags,
	Subcommands: []*cli.Command{
		activitiesCommand,
		athleteCommand,
		oauthCommand,
		rideCommand,
		uploadCommand,
	},
}

var AuthFlags = []cli.Flag{
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "cyclinganalytics.client-id",
		Usage: "API key for Cycling Analytics API",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "cyclinganalytics.client-secret",
		Usage: "API secret for Cycling Analytics API",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "cyclinganalytics.access-token",
		Usage: "Access token for Cycling Analytics API",
	}),
}
