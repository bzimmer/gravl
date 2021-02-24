package analysis

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/options"
)

type an struct {
	analyzer *analysis.Analyzer
	standard bool
}

var available = make(map[string]an)

// Add an analyzer; if standard is `true` it will be included in the standard run
// This function is not thread-safe; use it only on application startup
func Add(analyzer *analysis.Analyzer, standard bool) {
	if analyzer == nil {
		return
	}
	available[analyzer.Name] = an{analyzer: analyzer, standard: standard}
}

// All returns a copy of all the available analyzers
func All() []*analysis.Analyzer {
	var a []*analysis.Analyzer
	for _, val := range available {
		a = append(a, val.analyzer)
	}
	return a
}

func analyzers(c *cli.Context) ([]*analysis.Analyzer, error) {
	var ans []*analysis.Analyzer
	names := c.StringSlice("analyzer")
	if len(names) == 0 {
		for _, an := range available {
			if an.standard {
				ans = append(ans, an.analyzer)
			}
		}
		return ans, nil
	}
	for i := 0; i < len(names); i++ {
		opt, err := options.Parse(names[i])
		if err != nil {
			return nil, err
		}
		an, ok := available[opt.Name]
		if !ok {
			return nil, fmt.Errorf("unknown analyzer '%s'", opt.Name)
		}
		if an.analyzer.Flags != nil {
			if err := opt.ApplyFlags(an.analyzer.Flags); err != nil {
				return nil, err
			}
		}
		ans = append(ans, an.analyzer)
	}
	return ans, nil
}
