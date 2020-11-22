package gravl

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/bzimmer/httpwares"
	"github.com/markbates/goth"
	auth "github.com/markbates/goth/providers/strava"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/common/geo"
	"github.com/bzimmer/gravl/pkg/strava"
)

func newStravaAuthProvider(c *cli.Context, callback string) goth.Provider {
	provider := auth.New(
		c.String("strava.api-key"), c.String("strava.api-secret"), callback,
		// appears to be a bug where scope varargs do not work properly
		"read_all,profile:read_all,activity:read_all")
	t := http.DefaultTransport
	if c.Bool("http-tracing") {
		t = &httpwares.VerboseTransport{
			Transport: t,
		}
	}
	provider.HTTPClient = &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
	return provider
}

var stravaCommand = &cli.Command{
	Name:     "strava",
	Category: "route",
	Usage:    "Query Strava for rides and routes",
	Flags: []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "strava.api-key",
			Usage: "API key for Strava API",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "strava.api-secret",
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
	Action: func(c *cli.Context) error {
		client, err := strava.NewClient(
			strava.WithHTTPTracing(c.Bool("http-tracing")),
			strava.WithAPICredentials(
				c.String("strava.access-token"),
				c.String("strava.refresh-token")),
			strava.WithProvider(newStravaAuthProvider(c, "")))
		if err != nil {
			return err
		}

		if c.Bool("route") || c.Bool("activity") || c.Bool("stream") {
			args := c.Args()
			var tck geo.Trackable
			for i := 0; i < args.Len(); i++ {
				ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
				defer cancel()
				activityID, err := strconv.ParseInt(args.Get(i), 0, 64)
				if c.Bool("route") {
					tck, err = client.Route.Route(ctx, activityID)
				}
				if c.Bool("activity") {
					tck, err = client.Activity.Activity(ctx, activityID)
				}
				if c.Bool("stream") {
					tck, err = client.Activity.Streams(ctx, activityID, "latlng", "altitude")
				}
				if err != nil {
					return err
				}
				t, err := tck.Track()
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
			ath, err := client.Athlete.Athlete(ctx)
			if err != nil {
				return err
			}
			err = encoder.Encode(ath)
			if err != nil {
				return err
			}
			return nil
		}
		if c.Bool("refresh") {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			tokens, err := client.Auth.Refresh(ctx)
			if err != nil {
				return err
			}
			err = encoder.Encode(tokens)
			if err != nil {
				return err
			}
			return nil
		}
		if c.Bool("activities") {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			acts, err := client.Activity.Activities(ctx, strava.Pagination{Total: c.Int("count")})
			if err != nil {
				return err
			}
			for _, act := range acts {
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
			ath, err := client.Athlete.Athlete(c.Context)
			if err != nil {
				return err
			}
			rts, err := client.Route.Routes(ctx, ath.ID, strava.Pagination{Total: c.Int("count")})
			if err != nil {
				return err
			}
			for _, rt := range rts {
				err = encoder.Encode(rt)
				if err != nil {
					return err
				}
			}
			return nil
		}
		return nil
	},
}
