package rwgps

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/commands/activity/internal"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/activity"
	"github.com/bzimmer/gravl/pkg/providers/activity/rwgps"
)

func NewClient(c *cli.Context) (*rwgps.Client, error) {
	return rwgps.NewClient(
		rwgps.WithClientCredentials(c.String("rwgps.client-id"), ""),
		rwgps.WithTokenCredentials(c.String("rwgps.access-token"), "", time.Time{}),
		rwgps.WithHTTPTracing(c.Bool("http-tracing")),
	)
}

func athlete(c *cli.Context) error {
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	user, err := client.Users.AuthenticatedUser(ctx)
	if err != nil {
		return err
	}
	err = encoding.Encode(user)
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
	client, err := NewClient(c)
	if err != nil {
		return err
	}
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
	for _, trip := range trips {
		err = encoding.Encode(trip)
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
	client, err := NewClient(c)
	if err != nil {
		return err
	}
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
		if err = encoding.Encode(v); err != nil {
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

func poll(ctx context.Context, client *rwgps.Client, uploadID int64, follow bool) error {
	p := activity.NewPoller(client.Uploader())
	for res := range p.Poll(ctx, activity.UploadID(uploadID)) {
		if res.Err != nil {
			return res.Err
		}
		if err := encoding.Encode(res.Upload); err != nil {
			return err
		}
		if !follow {
			return nil
		}
	}
	return ctx.Err()
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
			u, err := client.Trips.Upload(ctx, file)
			if err != nil {
				return err
			}
			if !c.Bool("poll") {
				return encoding.Encode(u)
			}
			if err := poll(ctx, client, u.TaskID, true); err != nil {
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

var Command = &cli.Command{
	Name:     "rwgps",
	Category: "activity",
	Usage:    "Query RideWithGPS for rides and routes",
	Flags:    AuthFlags,
	Subcommands: []*cli.Command{
		activitiesCommand,
		activityCommand,
		athleteCommand,
		routeCommand,
		routesCommand,
		uploadCommand,
	},
}

var AuthFlags = []cli.Flag{
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "rwgps.client-id",
		Value: "",
		Usage: "Client ID for RideWithGPS API",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "rwgps.access-token",
		Value: "",
		Usage: "Access token for RideWithGPS API",
	}),
}
