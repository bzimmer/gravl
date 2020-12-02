package gravl

import (
	"context"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/common/geo"
	"github.com/bzimmer/gravl/pkg/strava"
)

var stravaCommand = &cli.Command{
	Name:     "strava",
	Category: "route",
	Usage:    "Query Strava for rides and routes",
	Flags:    stravaFlags,
	Action: func(c *cli.Context) error {
		client, err := strava.NewClient(
			strava.WithTokenCredentials(
				c.String("strava.access-token"),
				c.String("strava.refresh-token"),
				time.Now().Add(-1*time.Minute)),
			strava.WithClientCredentials(
				c.String("strava.client-id"),
				c.String("strava.client-secret")),
			strava.WithAutoRefresh(c.Context),
			strava.WithHTTPTracing(c.Bool("http-tracing")))
		if err != nil {
			return err
		}
		if c.Bool("route") || c.Bool("activity") || c.Bool("stream") {
			args := c.Args()
			var tracker geo.Tracker
			for i := 0; i < args.Len(); i++ {
				ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
				defer cancel()
				x, err := strconv.ParseInt(args.Get(i), 0, 64)
				if c.Bool("route") {
					tracker, err = client.Route.Route(ctx, x)
				} else if c.Bool("activity") {
					tracker, err = client.Activity.Activity(ctx, x)
				} else if c.Bool("stream") {
					tracker, err = client.Activity.Streams(ctx, x, "latlng", "altitude")
				}
				if err != nil {
					return err
				}
				t, err := tracker.Track()
				if err != nil {
					return err
				}
				if err = encoder.Encode(t); err != nil {
					return err
				}
			}
			return nil
		}
		if c.Bool("athlete") {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			athlete, err := client.Athlete.Athlete(ctx)
			if err != nil {
				return err
			}
			return encoder.Encode(athlete)
		}
		if c.Bool("refresh") {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			tokens, err := client.Auth.Refresh(ctx)
			if err != nil {
				return err
			}
			return encoder.Encode(tokens)
		}
		if c.Bool("activities") {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			activities, err := client.Activity.Activities(ctx, strava.Pagination{Total: c.Int("count")})
			if err != nil {
				return err
			}
			for _, act := range activities {
				err = encoder.Encode(act)
				if err != nil {
					return err
				}
			}
			return nil
		}
		if c.Bool("routes") {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			athlete, err := client.Athlete.Athlete(ctx)
			if err != nil {
				return err
			}
			routes, err := client.Route.Routes(ctx, athlete.ID, strava.Pagination{Total: c.Int("count")})
			if err != nil {
				return err
			}
			for _, route := range routes {
				err = encoder.Encode(route)
				if err != nil {
					return err
				}
			}
			return nil
		}
		return nil
	},
}

var stravaAuthFlags = []cli.Flag{
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "strava.client-id",
		Usage: "API key for Strava API",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "strava.client-secret",
		Usage: "API secret for Strava API",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "strava.access-token",
		Usage: "Access token for Strava API",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "strava.refresh-token",
		Usage: "Refresh token for Strava API",
	}),
}

var stravaFlags = merge(
	stravaAuthFlags,
	[]cli.Flag{
		&cli.BoolFlag{
			Name:    "activity",
			Aliases: []string{"t"},
			Value:   false,
			Usage:   "Activity",
		},
		&cli.BoolFlag{
			Name:    "stream",
			Aliases: []string{"s"},
			Value:   false,
			Usage:   "Stream",
		},
		&cli.BoolFlag{
			Name:    "route",
			Aliases: []string{"r"},
			Value:   false,
			Usage:   "Route",
		},
		&cli.BoolFlag{
			Name:    "athlete",
			Aliases: []string{"a"},
			Value:   false,
			Usage:   "Athlete",
		},
		&cli.BoolFlag{
			Name:    "activities",
			Aliases: []string{"A"},
			Value:   false,
			Usage:   "Activities",
		},
		&cli.BoolFlag{
			Name:    "routes",
			Aliases: []string{"R"},
			Value:   false,
			Usage:   "Routes",
		},
		&cli.IntFlag{
			Name:    "count",
			Aliases: []string{"N"},
			Value:   10,
			Usage:   "Count",
		},
		&cli.BoolFlag{
			Name:  "refresh",
			Value: false,
			Usage: "Refresh",
		},
	},
)
