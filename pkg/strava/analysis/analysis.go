package analysis

import (
	"context"
	"flag"
	"fmt"
)

type Analyzer struct {
	Name  string
	Doc   string
	Flags *flag.FlagSet
	Run   func(context.Context, *Pass) (interface{}, error)
}

func (a *Analyzer) String() string { return a.Name }

type Analysis struct {
	Args      []string
	Analyzers []*Analyzer
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

func (a *Analysis) Run(ctx context.Context, pass *Pass) (interface{}, error) {
	if err := a.applyFlags(); err != nil {
		return nil, err
	}

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
