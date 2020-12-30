package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/timshannon/bolthold"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"golang.org/x/time/rate"

	"github.com/bzimmer/gravl/pkg/common"
	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/store"
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
			strava.WithHTTPTracing(c.Bool("http-tracing")),
			strava.WithRateLimiter(
				rate.NewLimiter(rate.Every(1500*time.Millisecond), 25)),
			strava.WithCookieJar())
		if err != nil {
			return err
		}
		if c.Bool("update") {
			fn := c.Path("db")
			if fn == "" {
				return errors.New("nil db path")
			}
			directory := filepath.Dir(fn)
			if _, err := os.Stat(directory); os.IsNotExist(err) {
				log.Info().Str("directory", directory).Msg("creating")
				if err := os.MkdirAll(directory, os.ModeDir|0700); err != nil {
					return err
				}
			}
			db, err := bolthold.Open(fn, 0666, nil)
			if err != nil {
				return err
			}
			defer db.Close()
			log.Info().Str("db", fn).Msg("using database")

			var source store.Source
			if c.NArg() == 1 {
				source = &store.SourceFile{Path: c.Args().First()}
			} else {
				source = &store.SourceStrava{Client: client}
			}

			store := store.NewStore(db)
			n, err := store.Update(c.Context, source)
			if err != nil {
				return err
			}
			if err = encoder.Encode(map[string]int{"activities": n}); err != nil {
				return err
			}
			return nil
		}
		if c.Bool("route") || c.Bool("activity") || c.Bool("stream") {
			args := c.Args()
			var t interface{}
			for i := 0; i < args.Len(); i++ {
				ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
				defer cancel()
				x, err := strconv.ParseInt(args.Get(i), 0, 64)
				if c.Bool("route") {
					t, err = client.Route.Route(ctx, x)
				} else if c.Bool("activity") {
					t, err = client.Activity.Activity(ctx, x, "latlng", "altitude", "time")
				} else if c.Bool("stream") {
					t, err = client.Activity.Streams(ctx, x, "latlng", "altitude", "time")
				}
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
		if c.Bool("fitness") {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			username, password := c.String("strava.username"), c.String("strava.password")
			if err := client.Auth.Login(ctx, username, password); err != nil {
				return err
			}
			athlete, err := client.Athlete.Athlete(ctx)
			if err != nil {
				return err
			}
			foo, err := client.Fitness.TrainingLoad(ctx, athlete.ID)
			if err != nil {
				return err
			}
			return encoder.Encode(foo)
		}
		if c.Bool("activities") {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			defer func(start time.Time) {
				log.Debug().
					Dur("elapsed", time.Since(start)).
					Msg("activities")
			}(time.Now())
			activities, err := client.Activity.Activities(ctx, common.Pagination{Total: c.Int("count")})
			if err != nil {
				return err
			}
			for _, act := range activities {
				err = encoder.Encode(act)
				if err != nil {
					return err
				}
			}
		}
		if c.Bool("routes") {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			defer func(start time.Time) {
				log.Debug().
					Dur("elapsed", time.Since(start)).
					Msg("activities")
			}(time.Now())
			athlete, err := client.Athlete.Athlete(ctx)
			if err != nil {
				return err
			}
			routes, err := client.Route.Routes(ctx, athlete.ID, common.Pagination{Total: c.Int("count")})
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
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "strava.username",
		Usage: "Username for the Strava website",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "strava.password",
		Usage: "Password for the Strava website",
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
		&cli.BoolFlag{
			Name:    "fitness",
			Aliases: []string{"f"},
			Value:   false,
			Usage:   "Fitness profile",
		},
		&cli.IntFlag{
			Name:    "count",
			Aliases: []string{"N"},
			Value:   0,
			Usage:   "Count",
		},
		&cli.BoolFlag{
			Name:  "refresh",
			Value: false,
			Usage: "Refresh",
		},
		&cli.BoolFlag{
			Name:  "update",
			Value: false,
			Usage: "Update the databse with the latest activities",
		},
	},
)
