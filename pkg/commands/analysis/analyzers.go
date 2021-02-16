package analysis

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/analysis"
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
	"github.com/bzimmer/gravl/pkg/options"
)

type analyzer struct {
	analyzer *analysis.Analyzer
	standard bool
}

var available = func() map[string]analyzer {
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
