package gravl

import (
	"context"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/cyclinganalytics"
)

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
				time.Now().Add(-1*time.Minute)),
			cyclinganalytics.WithAutoRefresh(c.Context),
			cyclinganalytics.WithHTTPTracing(c.Bool("http-tracing")))
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		if c.Bool("athlete") {
			ath, merr := client.User.Me(ctx)
			if merr != nil {
				return merr
			}
			return encoder.Encode(ath)
		}
		if c.Bool("activities") {
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
	})}

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
	},
)
