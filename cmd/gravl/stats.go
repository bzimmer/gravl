package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/antonmedv/expr"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
	"github.com/urfave/cli/v2"
	"github.com/valyala/fastjson"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/analysis"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/climbing"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/eddington"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/festive500"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/hourrecord"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/koms"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/pythagorean"
	"github.com/bzimmer/gravl/pkg/strava/analysis/passes/splat"
)

type Env struct {
	Activities []*strava.Activity
}

var analyzers = []*analysis.Analyzer{
	climbing.New(),
	eddington.New(),
	festive500.New(),
	hourrecord.New(),
	koms.New(),
	pythagorean.New(),
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

func read(filename string) ([]*strava.Activity, error) {
	var (
		err   error
		sc    fastjson.Scanner
		acts  []*strava.Activity
		start = time.Now()
	)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	sc.InitBytes(b)
	for sc.Next() {
		if err = sc.Error(); err != nil {
			return nil, err
		}
		val := sc.Value()
		act := &strava.Activity{}
		err = json.Unmarshal(val.MarshalTo(nil), act)
		if err != nil {
			return nil, err
		}
		acts = append(acts, act)
	}
	log.Debug().
		Int("activities", len(acts)).
		Dur("elapsed", time.Since(start)).
		Msg("read")
	return acts, nil
}

// filter the activities
// For example:
//  {.Type in ["Ride"] && !.Commute && .StartDateLocal.Year() in [2020, 2019]}
func filter(c *cli.Context, acts []*strava.Activity) ([]*strava.Activity, error) {
	if !c.IsSet("filter") {
		return acts, nil
	}
	n := len(acts)
	start := time.Now()
	code := fmt.Sprintf("filter(Activities, %s)", closure(c.String("filter")))
	log.Debug().
		Str("code", code).
		Msg("filter")
	program, err := expr.Compile(code, expr.Env(Env{}))
	if err != nil {
		return nil, err
	}
	out, err := expr.Run(program, Env{Activities: acts})
	if err != nil {
		return nil, err
	}
	res := out.([]interface{})
	acts = make([]*strava.Activity, len(res))
	for i := range res {
		acts[i] = res[i].(*strava.Activity)
	}
	log.Debug().
		Int("activities{pre}", n).
		Int("activities{post}", len(acts)).
		Dur("elapsed", time.Since(start)).
		Msg("filter")
	return acts, nil
}

// groupby groups activities by a key
// currently only supports a single key for grouping
func groupby(c *cli.Context, acts []*strava.Activity) (map[string][]*strava.Activity, error) {
	if !c.IsSet("groupby") {
		return map[string][]*strava.Activity{
			"": acts,
		}, nil
	}
	start := time.Now()
	code := fmt.Sprintf("map(Activities, %s)", closure(c.String("groupby")))
	log.Debug().
		Str("code", code).
		Msg("groupby")
	program, err := expr.Compile(code, expr.Env(Env{}))
	if err != nil {
		return nil, err
	}
	out, err := expr.Run(program, Env{Activities: acts})
	if err != nil {
		return nil, err
	}
	res := out.([]interface{})
	groups := make(map[string][]*strava.Activity, len(res))
	for i, k := range res {
		key, err := cast.ToStringE(k)
		if err != nil {
			return nil, err
		}
		groups[key] = append(groups[key], acts[i])
	}
	log.Debug().
		Int("activities", len(acts)).
		Int("groups", len(groups)).
		Dur("elapsed", time.Since(start)).
		Msg("groupby")
	return groups, nil
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
		&cli.BoolFlag{
			Name:    "totals",
			Aliases: []string{"t"},
			Value:   false,
			Usage:   "Compute a total rather than grouped by years.",
		},
		&cli.StringSliceFlag{
			Name:    "analyzer",
			Aliases: []string{"a"},
			Usage:   "Analyzers to include (if none specified, default set is used)",
		},
	},
	Before: func(c *cli.Context) error {
		if c.NArg() == 0 {
			return errors.New("missing data file")
		}
		return nil
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
		any := analysis.Analysis{
			Args:      c.Args().Tail(),
			Analyzers: analyzers,
		}
		acts, err := read(c.Args().First())
		if err != nil {
			return err
		}
		acts, err = filter(c, acts)
		if err != nil {
			return err
		}
		groups, err := groupby(c, acts)
		if err != nil {
			return err
		}
		results := make(map[string]interface{})
		for key, group := range groups {
			pass := &analysis.Pass{Activities: group}
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
