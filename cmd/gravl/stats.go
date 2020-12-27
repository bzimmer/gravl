package main

import (
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/timshannon/bolthold"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/analysis"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/benford"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/climbing"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/eddington"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/festive500"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/hourrecord"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/kmeans"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/koms"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/pythagorean"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/splat"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/totals"
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
		hourrecord.New():  true,
		kmeans.New():      true,
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
	fn := c.Path("bolt")
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
	return &analysis.Pass{Activities: acts}, nil
}

// filter the activities
// For example:
//  {.Type in ["Ride"] && !.Commute && .StartDateLocal.Year() in [2020, 2019]}
func filter(c *cli.Context, pass *analysis.Pass) (*analysis.Pass, error) {
	if !c.IsSet("filter") {
		return pass, nil
	}
	q := closure(c.String("filter"))
	return pass.FilterExpr(q)
}

// groupby groups activities by a key
// currently only supports a single key for grouping
func groupby(c *cli.Context, pass *analysis.Pass) (map[string]*analysis.Pass, error) {
	if !c.IsSet("groupby") {
		return map[string]*analysis.Pass{
			"gravl": pass,
		}, nil
	}
	q := closure(c.String("groupby"))
	return pass.GroupByExpr(q)
}

var statsCommand = &cli.Command{
	Name:     "stats",
	Category: "route",
	Usage:    "Compute stats from Strava activities",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "filter",
			Aliases: []string{"f"},
			Usage:   "Expression for filtering activities",
		},
		&cli.StringFlag{
			Name:    "groupby",
			Aliases: []string{"g"},
			Usage:   "Expression for grouping activities",
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
		any := analysis.Analysis{
			Args:      c.Args().Tail(),
			Analyzers: as,
		}
		pass, err := read(c)
		if err != nil {
			return err
		}
		pass, err = filter(c, pass)
		if err != nil {
			return err
		}
		passes, err := groupby(c, pass)
		if err != nil {
			return err
		}
		results := make(map[string]interface{})
		for key, pass := range passes {
			res, err := any.Run(c.Context, pass)
			if err != nil {
				return err
			}
			results[key] = res
		}
		return encoder.Encode(results)
	},
}
