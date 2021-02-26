package strava

import (
	"context"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
	"golang.org/x/time/rate"

	"github.com/bzimmer/gravl/pkg/commands/activity"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	stravaweb "github.com/bzimmer/gravl/pkg/providers/activity/strava/web"
)

func NewWebClient(c *cli.Context) (*stravaweb.Client, error) {
	client, err := stravaweb.NewClient(
		stravaweb.WithHTTPTracing(c.Bool("http-tracing")),
		stravaweb.WithCookieJar(),
		stravaweb.WithRateLimiter(rate.NewLimiter(rate.Every(2*time.Second), 5)))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	username, password := c.String("strava.username"), c.String("strava.password")
	if err = client.Auth.Login(ctx, username, password); err != nil {
		return nil, err
	}
	return client, nil
}

func export(c *cli.Context) error {
	client, err := NewWebClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	args := c.Args().Slice()
	for i := 0; i < len(args); i++ {
		x, err := strconv.ParseInt(args[i], 0, 64)
		if err != nil {
			return err
		}
		exp, err := client.Export.Export(ctx, x)
		if err != nil {
			return err
		}
		if err = activity.Write(c, exp); err != nil {
			return err
		}
		if err = encoding.Encode(exp); err != nil {
			return err
		}
	}
	return nil
}

var exportCommand = &cli.Command{
	Name:      "export",
	Usage:     "Export a Strava activity by id",
	ArgsUsage: "ACTIVITY_ID (...)",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "overwrite",
			Aliases: []string{"o"},
			Value:   false,
			Usage:   "Overwrite the file if it exists; fail otherwise",
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"O"},
			Value:   "",
			Usage:   "The filename to use for writing the contents of the export, if not specified the contents are streamed to stdout",
		},
	},
	Action: export,
}

func fitness(c *cli.Context) error {
	webclient, err := NewWebClient(c)
	if err != nil {
		return err
	}
	apiclient, err := NewAPIClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	athlete, err := apiclient.Athlete.Athlete(ctx)
	if err != nil {
		return err
	}
	load, err := webclient.Fitness.TrainingLoad(ctx, athlete.ID)
	if err != nil {
		return err
	}
	return encoding.Encode(load)
}

var fitnessCommand = &cli.Command{
	Name:   "fitness",
	Usage:  "Query Strava for training load data for the authenticated user",
	Action: fitness,
}
