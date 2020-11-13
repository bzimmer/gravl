package gravl

import (
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/wta"
)

var wtaCommand = &cli.Command{
	Name:     "wta",
	Category: "geolocation",
	Usage:    "Query the WTA site for trip reports",
	Action: func(c *cli.Context) error {
		args := c.Args().Slice()
		if len(args) == 0 {
			// query the most recent if no reporter specified
			args = append(args, "")
		}
		client, err := wta.NewClient(
			wta.WithHTTPTracing(c.Bool("http-tracing")),
		)
		if err != nil {
			return err
		}
		for _, arg := range args {
			reports, err := client.Reports.TripReports(c.Context, arg)
			if err != nil {
				return err
			}
			err = encoder.Encode(reports)
			if err != nil {
				return err
			}
		}
		return nil
	},
}
