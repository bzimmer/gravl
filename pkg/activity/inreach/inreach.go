package inreach

import (
	"errors"
	"sync"
	"time"

	"github.com/bzimmer/activity/inreach"
	"github.com/rs/zerolog/log"
	"github.com/tj/go-naturaldate"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg"
)

const Provider = "inreach"

var before sync.Once

func since(c *cli.Context) (inreach.APIOption, error) {
	var opt inreach.APIOption
	if c.IsSet("since") {
		before := time.Now()
		after, err := naturaldate.Parse(c.String("since"), before)
		if err != nil {
			return nil, err
		}
		log.Info().Time("before", before).Time("after", after).Msg("date range")
		if after.After(before) {
			return nil, errors.New("invalid date range")
		}
		opt = inreach.WithDateRange(before, after)
	}
	return opt, nil
}

func activities(c *cli.Context) error {
	opt, err := since(c)
	if err != nil {
		return err
	}
	enc := pkg.Runtime(c).Encoder
	met := pkg.Runtime(c).Metrics
	client := pkg.Runtime(c).InReach
	for i := 0; i < c.NArg(); i++ {
		arg := c.Args().Get(i)
		feed, err := client.Feed.Feed(c.Context, arg, opt)
		if err != nil {
			return err
		}
		collection, err := feed.GeoJSON()
		if err != nil {
			return err
		}
		met.IncrCounter([]string{Provider, c.Command.Name}, 1)
		if err := enc.Encode(collection); err != nil {
			return err
		}
	}
	return nil
}

func activitiesCommand() *cli.Command {
	return &cli.Command{
		Name:    "activities",
		Usage:   "Query activities for a user from InReach",
		Aliases: []string{"A"},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "count",
				Aliases: []string{"N"},
				Value:   0,
				Usage:   "The number of activities to query from InReach (the number returned will be <= N)",
			},
			&cli.StringFlag{
				Name:  "since",
				Usage: "Return results since the time specified",
			},
		},
		Action: activities,
	}
}

// Before configures the InReach client
func Before(c *cli.Context) error {
	var err error
	before.Do(func() {
		var client *inreach.Client
		client, err = inreach.NewClient(
			inreach.WithHTTPTracing(c.Bool("http-tracing")))
		if err != nil {
			return
		}
		pkg.Runtime(c).InReach = client
		log.Info().Msg("created inreach client")
	})
	return err
}

func Command() *cli.Command {
	return &cli.Command{
		Name:        "inreach",
		Category:    "activity",
		Usage:       "Query InReach for activities",
		Description: "Operations supported by the InReach KML API",
		Before:      Before,
		Subcommands: []*cli.Command{
			activitiesCommand(),
		},
	}
}
