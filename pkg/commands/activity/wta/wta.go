package wta

import (
	"context"

	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/activity/wta"
)

var Command = &cli.Command{
	Name:      "wta",
	Category:  "activity",
	Usage:     "Query the WTA site for trip reports, if no reporter is specified the most recent reports are returned",
	ArgsUsage: "[REPORTER_NAME ...]",
	Action: func(c *cli.Context) error {
		args := c.Args().Slice()
		if len(args) == 0 {
			// query the most recent if no reporter specified
			args = append(args, "")
		}
		client, err := wta.NewClient(wta.WithHTTPTracing(c.Bool("http-tracing")))
		if err != nil {
			return err
		}
		for _, arg := range args {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			reports, err := client.Reports.TripReports(ctx, arg)
			if err != nil {
				return err
			}
			if err = encoding.Encode(reports); err != nil {
				return err
			}
		}
		return nil
	},
}
