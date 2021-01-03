package main

import (
	"context"
	"errors"
	"math"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/timshannon/bolthold"
	"github.com/urfave/cli/v2"

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
	"github.com/bzimmer/gravl/pkg/analysis/passes/totals"
	"github.com/bzimmer/gravl/pkg/strava"
)

type analyzer struct {
	analyzer *analysis.Analyzer
	standard bool
}

var analyzers = func() map[string]analyzer {
	res := make(map[string]analyzer)
	for an, standard := range map[*analysis.Analyzer]bool{
		benford.New():     false,
		climbing.New():    true,
		eddington.New():   true,
		festive500.New():  true,
		forecast.New():    false,
		hourrecord.New():  true,
		cluster.New():     false,
		koms.New():        true,
		pythagorean.New(): true,
		splat.New():       false,
		totals.New():      true,
	} {
		res[an.Name] = analyzer{analyzer: an, standard: standard}
	}
	return res
}()

func closure(f string) string {
	if f == "" {
		return f
	}
	if !strings.HasPrefix(f, "{") {
		f = "{" + f
	}
	if !strings.HasSuffix(f, "}") {
		f = f + "}"
	}
	return f
}

func read(c *cli.Context) (*analysis.Pass, error) {
	fn := c.Path("db")
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
	q := closure(c.String("filter"))
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
		exprs = append(exprs, closure(g))
	}
	g, err := pass.GroupBy(exprs...)
	if err != nil {
		return nil, err
	}
	return g, nil
}

var statsCommand = &cli.Command{
	Name:     "stats",
	Category: "route",
	Usage:    "Compute stats from Strava activities",
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
	},
	Action: func(c *cli.Context) error {
		var as []*analysis.Analyzer
		if c.IsSet("analyzer") {
			names := c.StringSlice("analyzer")
			for i := 0; i < len(names); i++ {
				an, ok := analyzers[names[i]]
				if !ok {
					log.Warn().Str("name", names[i]).Msg("missing analyzer")
					continue
				}
				as = append(as, an.analyzer)
			}
		} else {
			for _, an := range analyzers {
				if an.standard {
					as = append(as, an.analyzer)
				}
			}
		}
		if len(as) == 0 {
			return errors.New("no analyzers found")
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
		w := &Analyzer{
			Any: analysis.Analysis{
				Args:      c.Args().Tail(),
				Analyzers: as,
			}}
		if err = group.Walk(c.Context, w.Run); err != nil {
			return err
		}
		return encoder.Encode(w.Collect())
	},
}

type Results struct {
	Key     string
	Level   int
	Results interface{}
}

type Analyzer struct {
	Any     analysis.Analysis
	Results []*Results
}

func (a *Analyzer) Run(ctx context.Context, g *analysis.Group) error {
	if len(g.Groups) > 0 {
		// not a leaf node so skip evaluation
		a.Results = append(a.Results, &Results{Key: g.Key, Level: g.Level})
		return nil
	}
	res, err := a.Any.Run(ctx, g.Pass)
	if err != nil {
		return err
	}
	a.Results = append(a.Results, &Results{Key: g.Key, Results: res, Level: g.Level})
	return nil
}

func (a *Analyzer) Collect() map[string]interface{} {
	var res []map[string]interface{}
	for _, x := range a.Results {
		for len(res) > x.Level {
			res = res[:len(res)-1]
		}
		if len(res) == x.Level {
			m := make(map[string]interface{})
			if len(res) > 0 {
				res[len(res)-1][x.Key] = m
			}
			res = append(res, m)
		}
		if x.Results != nil {
			n := int(math.Max(float64(x.Level-1), 0))
			res[n][x.Key] = x.Results
		}
	}
	if len(res) == 0 {
		return nil
	}
	return res[0]
}
