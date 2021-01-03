package analysis

import (
	"context"
	"flag"
	"fmt"
	"math"
)

type Analyzer struct {
	Name  string
	Doc   string
	Flags *flag.FlagSet
	Run   func(context.Context, *Pass) (interface{}, error)
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

	analyzer := ""
	flags := make(map[string][]string)
	for i := 0; i < len(a.Args); i++ {
		arg := a.Args[i]
		if arg == "--" {
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

func (a *Analysis) RunPass(ctx context.Context, pass *Pass) (interface{}, error) {
	results := make(map[string]interface{})
	for _, analyzer := range a.Analyzers {
		res, err := analyzer.Run(ctx, pass)
		if err != nil {
			return nil, err
		}
		results[analyzer.Name] = res
	}
	return results, nil
}

func (a *Analysis) RunGroup(ctx context.Context, group *Group) (map[string]interface{}, error) {
	if err := group.Walk(ctx, a.runGroup); err != nil {
		return nil, err
	}
	return a.collect(), nil
}

func (a *Analysis) runGroup(ctx context.Context, g *Group) error {
	if len(g.Groups) > 0 {
		// not a leaf node so skip evaluation
		a.results = append(a.results, &results{Key: g.Key, Level: g.Level})
		return nil
	}
	res, err := a.RunPass(ctx, g.Pass)
	if err != nil {
		return err
	}
	a.results = append(a.results, &results{Key: g.Key, Results: res, Level: g.Level})
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
