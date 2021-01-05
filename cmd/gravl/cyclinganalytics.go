package main

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/activity/cyclinganalytics"
)

// collect returns a slice of files for uploading
// Primary use case has been uploading fit files from Zwift so this function
//  filters small files (less then 1K) and files of the name "inProgressActivity.fit"
func collect(name string) ([]*cyclinganalytics.File, error) {
	var files []*cyclinganalytics.File
	err := filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			base := filepath.Base(path)
			if base == "inProgressActivity.fit" {
				log.Warn().
					Str("name", path).
					Msg("skipping, not a completed activity")
				return nil
			}
			if info.Size() < 1024 {
				log.Warn().
					Int64("size", info.Size()).
					Str("name", path).
					Msg("skipping, too small")
				return nil
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			files = append(files, &cyclinganalytics.File{
				Name:   base,
				Reader: file,
				Size:   info.Size(),
			})
		}
		return nil
	})
	return files, err
}

var cyclinganalyticsCommand = &cli.Command{
	Name:     "cyclinganalytics",
	Aliases:  []string{"ca"},
	Category: "route",
	Usage:    "Query the cyclinganalytics.com site",
	Flags:    cyclingAnalyticsFlags,
	Action: func(c *cli.Context) error {
		client, err := cyclinganalytics.NewClient(
			cyclinganalytics.WithTokenCredentials(
				c.String("cyclinganalytics.access-token"),
				c.String("cyclinganalytics.refresh-token"),
				time.Time{}),
			cyclinganalytics.WithAutoRefresh(c.Context),
			cyclinganalytics.WithHTTPTracing(c.Bool("http-tracing")))
		if err != nil {
			return err
		}

		if c.Bool("upload") {
			args := c.Args()
			for i := 0; i < args.Len(); i++ {
				files, err := collect(args.Get(i))
				if err != nil {
					return err
				}
				for _, file := range files {
					defer file.Close()
					ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
					defer cancel()
					log.Info().
						Str("file", file.Name).
						Int64("size", file.Size).
						Msg("uploading")
					u, err := client.Rides.Upload(ctx, cyclinganalytics.Me, file)
					if err != nil {
						return err
					}
					if err := encoder.Encode(u); err != nil {
						return err
					}
				}
			}
			return nil
		}
		if c.Bool("status") {
			args := c.Args()
			for i := 0; i < args.Len(); i++ {
				ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
				defer cancel()
				uploadID, err := strconv.ParseInt(args.Get(i), 0, 64)
				if err != nil {
					return err
				}
				status, err := client.Rides.Status(ctx, cyclinganalytics.Me, uploadID)
				if err != nil {
					return err
				}
				if err = encoder.Encode(status); err != nil {
					return err
				}
			}
		}
		if c.Bool("athlete") {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			athlete, err := client.User.Me(ctx)
			if err != nil {
				return err
			}
			return encoder.Encode(athlete)
		}
		if c.Bool("activities") {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			rides, err := client.Rides.Rides(ctx, cyclinganalytics.Me)
			if err != nil {
				return err
			}
			for _, ride := range rides {
				err := encoder.Encode(ride)
				if err != nil {
					return err
				}
			}
		}
		if c.Bool("activity") {
			args := c.Args()
			opts := cyclinganalytics.RideOptions{
				Streams: []string{"latitude", "longitude", "elevation"},
			}
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
				if err = encoder.Encode(ride); err != nil {
					return err
				}
			}
		}
		return nil
	},
}

var cyclingAnalyticsAuthFlags = []cli.Flag{
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

var cyclingAnalyticsFlags = merge(
	cyclingAnalyticsAuthFlags,
	[]cli.Flag{
		&cli.BoolFlag{
			Name:    "athlete",
			Aliases: []string{"a"},
			Value:   false,
			Usage:   "Athlete",
		},
		&cli.BoolFlag{
			Name:    "activity",
			Aliases: []string{"t"},
			Value:   false,
			Usage:   "Activity",
		},
		&cli.BoolFlag{
			Name:    "activities",
			Aliases: []string{"A"},
			Value:   false,
			Usage:   "Activities",
		},
		&cli.BoolFlag{
			Name:    "upload",
			Aliases: []string{"u"},
			Value:   false,
			Usage:   "Upload",
		},
		&cli.BoolFlag{
			Name:    "status",
			Aliases: []string{"U"},
			Value:   false,
			Usage:   "Upload status",
		},
	},
)
