package gravl

import (
	"context"
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
			cyclinganalytics.WithHTTPTracing(c.Bool("http-tracing")),
			cyclinganalytics.WithTokenCredentials(
				c.String("cyclinganalytics.access-token"),
				c.String("cyclinganalytics.refresh-token"),
				time.Time{}))
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		rides, err := client.Rides.Rides(ctx)
		if err != nil {
			return err
		}
		for _, ride := range rides {
			err := encoder.Encode(ride)
			if err != nil {
				return err
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
	cyclingAnalyticsAuthFlags)
