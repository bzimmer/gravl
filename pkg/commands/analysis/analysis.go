package analysis

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/timshannon/bolthold"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/activity/strava"
	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/analysis/passes/benford"
	"github.com/bzimmer/gravl/pkg/analysis/passes/climbing"
	"github.com/bzimmer/gravl/pkg/analysis/passes/cluster"
	"github.com/bzimmer/gravl/pkg/analysis/passes/eddington"
	"github.com/bzimmer/gravl/pkg/analysis/passes/festive500"
	"github.com/bzimmer/gravl/pkg/analysis/passes/forecast"
	"github.com/bzimmer/gravl/pkg/analysis/passes/hourrecord"
	"github.com/bzimmer/gravl/pkg/analysis/passes/koms"
	"github.com/bzimmer/gravl/pkg/analysis/passes/pythagorean"
	"github.com/bzimmer/gravl/pkg/analysis/passes/splat"
	"github.com/bzimmer/gravl/pkg/analysis/passes/staticmap"
	"github.com/bzimmer/gravl/pkg/analysis/passes/totals"
	"github.com/bzimmer/gravl/pkg/commands"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
)

type analyzer struct {
	analyzer *analysis.Analyzer
	standard bool
}

var _analyzers = func() map[string]analyzer {
	res := make(map[string]analyzer)
	for an, standard := range map[*analysis.Analyzer]bool{
		benford.New():     false,
		climbing.New():    true,
		cluster.New():     false,
		eddington.New():   true,
		festive500.New():  true,
		forecast.New():    false,
		hourrecord.New():  true,
		koms.New():        true,
		pythagorean.New(): true,
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

func read(c *cli.Context) (*analysis.Pass, error) {
	fn := c.Path("store")
	if fn == "" {
		return nil, errors.New("nil db path")
	}
	store, err := bolthold.Open(fn, 0666, nil)
	if err != nil {
		return nil, err
	}
	defer store.Close()

	var acts []*strava.Activity
	err = store.ForEach(&bolthold.Query{}, func(act *strava.Activity) error {
		acts = append(acts, act)
		return nil
	})
	if err != nil {
		return nil, err
	}
	uf := c.Generic("units").(*analysis.UnitsFlag)
	return &analysis.Pass{Activities: acts, Units: uf.Units}, nil
}

// filter the activities
// For example:
//  {.Type in ["Ride"] && !.Commute && .StartDateLocal.Year() in [2020, 2019]}
func filter(c *cli.Context, pass *analysis.Pass) (*analysis.Pass, error) {
	if !c.IsSet("filter") {
		return pass, nil
	}
	q := analysis.Closure(c.String("filter"))
	return pass.Filter(q)
}

// groupby groups activities by expression values
func groupby(c *cli.Context, pass *analysis.Pass) (*analysis.Group, error) {
	if !c.IsSet("groupby") {
		return &analysis.Group{
			Pass: pass,
		}, nil
	}
	var exprs []string
	for _, g := range c.StringSlice("groupby") {
		exprs = append(exprs, analysis.Closure(g))
	}
	g, err := pass.GroupBy(exprs...)
	if err != nil {
		return nil, err
	}
	return g, nil
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
			Name:    "groupby",
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
		ans, err := analyzers(c)
		if err != nil {
			return err
		}
		pass, err := read(c)
		if err != nil {
			return err
		}
		pass, err = filter(c, pass)
		if err != nil {
			return err
		}
		group, err := groupby(c, pass)
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
		results, err := any.RunGroup(ctx, group)
		if err != nil {
			return err
		}
		return encoding.Encode(results)
	},
}
