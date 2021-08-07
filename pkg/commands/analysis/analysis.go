package analysis

import (
	"context"

	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/commands"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/commands/store"
	"github.com/bzimmer/gravl/pkg/eval"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

// read activities from the store
func read(ctx context.Context, acts <-chan *strava.ActivityResult) ([]*strava.Activity, error) {
	var activities []*strava.Activity
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case res, ok := <-acts:
			if !ok {
				return activities, nil
			}
			if res.Err != nil {
				return nil, res.Err
			}
			if res.Activity == nil {
				continue
			}
			activities = append(activities, res.Activity)
		}
	}
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
	filterer, err := commands.Filterer(c.String("filter"))
	if err != nil {
		return nil, err
	}
	return filterer.Filter(c.Context, acts)
}

// group groups activities by expression values
//
// The result of the expression will be converted a string and used as the key
// in the final result map.
func group(c *cli.Context, acts []*strava.Activity) (*analysis.Pass, error) {
	var mappers []eval.Mapper
	for _, q := range c.StringSlice("group") {
		mapper, err := commands.Mapper(q)
		if err != nil {
			return nil, err
		}
		mappers = append(mappers, mapper)
	}
	return analysis.Group(c.Context, acts, mappers...)
}

func analyze(c *cli.Context) error {
	ctx, cancel := context.WithCancel(c.Context)
	defer cancel()
	db, err := store.Open(c, "input")
	if err != nil {
		return err
	}
	acts, err := read(ctx, db.Activities(ctx))
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
	if c.IsSet("timeout") {
		ctx, cancel = context.WithTimeout(ctx, c.Duration("timeout"))
		defer cancel()
	}
	uf := c.Generic("units").(*analysis.UnitsFlag)
	x := analysis.WithContext(ctx, uf.Units)

	any := analysis.NewAnalysis(ans)
	results, err := any.Run(x, pass)
	if err != nil {
		return err
	}
	return encoding.For(c).Encode(results)
}

var listCommand = &cli.Command{
	Name:  "list",
	Usage: "Return the list of available analyzers",
	Action: func(c *cli.Context) error {
		res := make(map[string]map[string]interface{})
		for nm, an := range available {
			res[nm] = make(map[string]interface{})
			res[nm]["doc"] = an.analyzer.Doc
			res[nm]["base"] = an.standard
			res[nm]["flags"] = (an.analyzer.Flags != nil)
		}
		return encoding.For(c).Encode(res)
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
		}, &cli.StringFlag{
			Name:    "filter",
			Aliases: []string{"f"},
			Usage:   "Expression for filtering activities to remove",
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
		store.InputFlag(store.DefaultLocalStore),
	},
	Subcommands: []*cli.Command{listCommand},
	Action:      analyze,
}
