package analysis

import (
	"context"
	"errors"

	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/analysis/eval"
	"github.com/bzimmer/gravl/pkg/analysis/eval/antonmedv"
	"github.com/bzimmer/gravl/pkg/analysis/store/bunt"
	"github.com/bzimmer/gravl/pkg/commands"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

func read(c *cli.Context) ([]*strava.Activity, error) {
	path := c.Path("store")
	if path == "" {
		return nil, errors.New("nil db path")
	}
	db, err := bunt.Open(path)
	if err != nil {
		return nil, err
	}
	ca, ce := db.Activities(c.Context)
	acts, err := strava.Activities(c.Context, ca, ce)
	if err != nil {
		return nil, err
	}
	return acts, nil
}

// filter the activities
//
// The expression must return a boolean value.
//
// For example:
//  .Type in ["Ride"] && !.Commute && .StartDateLocal.Year() in [2020, 2019]
func filter(c *cli.Context, acts []*strava.Activity) ([]*strava.Activity, error) {
	if !c.IsSet("filter") {
		return acts, nil
	}
	evaluator := antonmedv.New(c.String("filter"))
	return evaluator.Filter(c.Context, acts)
}

// group groups activities by expression values
//
// The result of the expression will be converted a string and used as the key
// in the final result map.
func group(c *cli.Context, acts []*strava.Activity) (*analysis.Pass, error) {
	var expressions []eval.Mapper
	for _, q := range c.StringSlice("group") {
		expressions = append(expressions, antonmedv.New(q))
	}
	return analysis.Group(c.Context, acts, expressions...)
}

var Command = &cli.Command{
	Name:     "analysis",
	Aliases:  []string{"pass"},
	Category: "analysis",
	Usage:    "Produce statistics and other interesting artifacts from Strava activities",
	Flags: []cli.Flag{
		&cli.GenericFlag{
			Name:    "units",
			Aliases: []string{"u"},
			Usage:   "Units",
			Value:   &analysis.UnitsFlag{},
		},
		&cli.StringFlag{
			Name:    "filter",
			Aliases: []string{"f"},
			Usage:   "Expression for filtering activities",
		},
		&cli.StringSliceFlag{
			Name:    "group",
			Aliases: []string{"g"},
			Usage:   "Expressions for grouping activities",
		},
		&cli.StringSliceFlag{
			Name:    "analyzer",
			Aliases: []string{"a"},
			Usage:   "Analyzers to include (if none specified, default set is used)",
		},
		commands.StoreFlag,
	},
	Action: func(c *cli.Context) error {
		acts, err := read(c)
		if err != nil {
			return err
		}
		acts, err = filter(c, acts)
		if err != nil {
			return err
		}
		pass, err := group(c, acts)
		if err != nil {
			return err
		}
		ans, err := analyzers(c)
		if err != nil {
			return err
		}
		any, err := analysis.NewAnalysis(ans, c.Args().Slice())
		if err != nil {
			return err
		}
		ctx := c.Context
		if c.IsSet("timeout") {
			x, cancel := context.WithTimeout(ctx, c.Duration("timeout"))
			defer cancel()
			ctx = x
		}
		uf := c.Generic("units").(*analysis.UnitsFlag)
		x := analysis.WithContext(ctx, uf.Units)
		results, err := any.Run(x, pass)
		if err != nil {
			return err
		}
		return encoding.Encode(results)
	},
}
