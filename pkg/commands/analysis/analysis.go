package analysis

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/analysis/eval"
	"github.com/bzimmer/gravl/pkg/analysis/eval/antonmedv"
	"github.com/bzimmer/gravl/pkg/analysis/passes/ageride"
	"github.com/bzimmer/gravl/pkg/analysis/passes/benford"
	"github.com/bzimmer/gravl/pkg/analysis/passes/climbing"
	"github.com/bzimmer/gravl/pkg/analysis/passes/cluster"
	"github.com/bzimmer/gravl/pkg/analysis/passes/eddington"
	"github.com/bzimmer/gravl/pkg/analysis/passes/festive500"
	"github.com/bzimmer/gravl/pkg/analysis/passes/forecast"
	"github.com/bzimmer/gravl/pkg/analysis/passes/hourrecord"
	"github.com/bzimmer/gravl/pkg/analysis/passes/koms"
	"github.com/bzimmer/gravl/pkg/analysis/passes/pythagorean"
	"github.com/bzimmer/gravl/pkg/analysis/passes/rolling"
	"github.com/bzimmer/gravl/pkg/analysis/passes/splat"
	"github.com/bzimmer/gravl/pkg/analysis/passes/staticmap"
	"github.com/bzimmer/gravl/pkg/analysis/passes/totals"
	"github.com/bzimmer/gravl/pkg/analysis/store/bunt"
	"github.com/bzimmer/gravl/pkg/commands"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

type analyzer struct {
	analyzer *analysis.Analyzer
	standard bool
}

var _analyzers = func() map[string]analyzer {
	res := make(map[string]analyzer)
	for an, standard := range map[*analysis.Analyzer]bool{
		ageride.New():     false,
		benford.New():     false,
		climbing.New():    true,
		cluster.New():     false,
		eddington.New():   true,
		festive500.New():  true,
		forecast.New():    false,
		hourrecord.New():  true,
		koms.New():        true,
		pythagorean.New(): true,
		rolling.New():     true,
		splat.New():       false,
		staticmap.New():   false,
		totals.New():      true,
	} {
		res[an.Name] = analyzer{analyzer: an, standard: standard}
	}
	return res
}()

func analyzers(c *cli.Context) ([]*analysis.Analyzer, error) {
	var ans []*analysis.Analyzer
	if c.IsSet("analyzer") {
		names := c.StringSlice("analyzer")
		for i := 0; i < len(names); i++ {
			an, ok := _analyzers[names[i]]
			if !ok {
				log.Warn().Str("name", names[i]).Msg("missing analyzer")
				continue
			}
			ans = append(ans, an.analyzer)
		}
	} else {
		for _, an := range _analyzers {
			if an.standard {
				ans = append(ans, an.analyzer)
			}
		}
	}
	if len(ans) == 0 {
		return nil, errors.New("no analyzers found")
	}
	return ans, nil
}

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
func group(c *cli.Context, acts []*strava.Activity) (*analysis.Pass, error) {
	var expressions []eval.Evaluator
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
