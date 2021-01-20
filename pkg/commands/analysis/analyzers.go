package analysis

import (
	"errors"

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
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type analyzer struct {
	analyzer *analysis.Analyzer
	standard bool
}

var _analyzers = func() map[string]analyzer {
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
