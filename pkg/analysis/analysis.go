package analysis

import (
	"flag"
	"fmt"
	"math"

	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

type Analyzer struct {
	Name  string
	Doc   string
	Flags *flag.FlagSet
	Run   func(*Context, []*strava.Activity) (interface{}, error)
}

func (a *Analyzer) String() string { return a.Name }

type results struct {
	Key     string
	Level   int
	Results interface{}
}

type Analysis struct {
	Args      []string
	Analyzers []*Analyzer
	results   []*results
}

func NewAnalysis(analyzers []*Analyzer, args []string) (*Analysis, error) {
	a := &Analysis{
		Args:      args,
		Analyzers: analyzers,
	}
	if err := a.applyFlags(); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Analysis) applyFlags() error {
	if len(a.Args) == 0 {
		return nil
	}

	analyzers := make(map[string]*Analyzer)
	for _, y := range a.Analyzers {
		if y.Flags == nil {
			continue
		}
		analyzers[y.Name] = y
	}

	dashes, analyzer := false, ""
	flags := make(map[string][]string)
	for i := 0; i < len(a.Args); i++ {
		arg := a.Args[i]
		if arg == "--" {
			dashes = true
			continue
		}
		if !dashes {
			continue
		}
		if analyzer == "" {
			// this arg should be an analyzer name
			_, ok := analyzers[arg]
			if !ok {
				return fmt.Errorf("expected analyzer name, found '%s'", arg)
			}
			analyzer = arg
			continue
		}
		_, ok := analyzers[arg]
		if ok {
			// starts a set of flags for this analyzer
			analyzer = arg
			continue
		}
		flags[analyzer] = append(flags[analyzer], arg)
	}

	// apply the flags to the analyzer
	for key, values := range flags {
		if err := analyzers[key].Flags.Parse(values); err != nil {
			return err
		}
	}

	return nil
}

func (a *Analysis) Run(ctx *Context, pass *Pass) (map[string]interface{}, error) {
	if err := a.run(ctx, pass, 0); err != nil {
		return nil, err
	}
	return a.collect(), nil
}

func (a *Analysis) run(ctx *Context, pass *Pass, level int) error {
	if len(pass.Children) > 0 {
		a.results = append(a.results, &results{Key: pass.Key, Level: level})
		for _, child := range pass.Children {
			if err := a.run(ctx, child, level+1); err != nil {
				return err
			}
		}
		return nil
	}
	res := make(map[string]interface{})
	for _, analyzer := range a.Analyzers {
		r, err := analyzer.Run(ctx, pass.Activities)
		if err != nil {
			return err
		}
		res[analyzer.Name] = r
	}
	a.results = append(a.results, &results{Key: pass.Key, Results: res, Level: level})
	return nil
}

func (a *Analysis) collect() map[string]interface{} {
	var res []map[string]interface{}
	for _, x := range a.results {
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
