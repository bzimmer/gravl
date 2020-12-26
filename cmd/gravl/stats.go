package main

import (
	"errors"
	"strings"

	bh "github.com/timshannon/bolthold"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/analysis"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/climbing"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/eddington"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/festive500"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/hourrecord"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/kmeans"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/koms"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/pythagorean"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/splat"
)

var analyzers = []*analysis.Analyzer{
	climbing.New(),
	eddington.New(),
	festive500.New(),
	hourrecord.New(),
	koms.New(),
	pythagorean.New(),
	kmeans.New(),
}

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
	store, err := bh.Open(fn, 0666, nil)
	if err != nil {
		return nil, err
	}
	defer store.Close()

	var acts []*strava.Activity
	err = store.ForEach(&bh.Query{}, func(act *strava.Activity) error {
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
			"": pass,
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
		if c.IsSet("analyzer") {
			any := splat.New()
			names := c.StringSlice("analyzer")
			var anys []*analysis.Analyzer
			for i := 0; i < len(names); i++ {
				if names[i] == any.Name {
					anys = append(anys, any)
					continue
				}
				for j := 0; j < len(analyzers); j++ {
					if names[i] == analyzers[j].Name {
						anys = append(anys, analyzers[j])
					}
				}
			}
			analyzers = anys
		}
		if len(analyzers) == 0 {
			return errors.New("no analyzers configured")
		}
		any := analysis.Analysis{
			Args:      c.Args().Tail(),
			Analyzers: analyzers,
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
		// special case, if one group and the key is `""`, return as a list not a map
		if val, ok := results[""]; ok && len(results) == 1 {
			return encoder.Encode(val)
		}
		return encoder.Encode(results)
	},
}
