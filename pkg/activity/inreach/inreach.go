package inreach

import (
	"errors"
	"sync"

	"github.com/bzimmer/activity/inreach"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/activity"
)

const Provider = "inreach"

var before sync.Once

func daterange(c *cli.Context) (inreach.APIOption, error) {
	before, after, err := activity.DateRange(c)
	if err != nil {
		return nil, err
	}
	log.Info().Time("before", before).Time("after", after).Msg("date range")
	return inreach.WithDateRange(before, after), nil
}

func activities(c *cli.Context) error {
	opt, err := daterange(c)
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
		Flags:   activity.DateRangeFlags(),
		Before: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return errors.New("no user specified")
			}
			return nil
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

// Command returns a fully instantiated cli command
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
