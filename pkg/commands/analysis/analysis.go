package analysis

import (
	"context"

	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/analysis/eval"
	"github.com/bzimmer/gravl/pkg/analysis/store"
	"github.com/bzimmer/gravl/pkg/commands"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	storecmd "github.com/bzimmer/gravl/pkg/commands/store"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

func read(c *cli.Context, db store.Store) ([]*strava.Activity, error) {
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
	filterer := commands.Filterer(c.String("filter"))
	return filterer.Filter(c.Context, acts)
}

// group groups activities by expression values
//
// The result of the expression will be converted a string and used as the key
// in the final result map.
func group(c *cli.Context, acts []*strava.Activity) (*analysis.Pass, error) {
	var mappers []eval.Mapper
	for _, q := range c.StringSlice("group") {
		mappers = append(mappers, commands.Mapper(q))
	}
	return analysis.Group(c.Context, acts, mappers...)
}

func analyze(c *cli.Context) error {
	db, err := storecmd.Open(c, "input", storecmd.DefaultLocalStore)
	if err != nil {
		return err
	}
	acts, err := read(c, db)
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
}

var listCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{""},
	Usage:   "Return the list of available analyzers",
	Action: func(c *cli.Context) error {
		res := make(map[string]map[string]interface{})
		for nm, an := range _analyzers {
			res[nm] = make(map[string]interface{})
			res[nm]["doc"] = an.analyzer.Doc
			res[nm]["base"] = an.standard
			res[nm]["flags"] = (an.analyzer.Flags != nil)
		}
		return encoding.Encode(res)
	},
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
		storecmd.InputFlag(storecmd.DefaultLocalStore),
	},
	Subcommands: []*cli.Command{listCommand},
	Action:      analyze,
}
