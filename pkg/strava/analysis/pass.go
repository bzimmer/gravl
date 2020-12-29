package analysis

import (
	"fmt"
	"time"

	"github.com/antonmedv/expr"
	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
)

type Pass struct {
	Units      Units
	Activities []*strava.Activity
}

// FilterExpr filters the activities using an expression
// For example:
//  {.Type in ["Ride"] && !.Commute && .StartDateLocal.Year() in [2020, 2019]}
func (p *Pass) FilterExpr(q string) (*Pass, error) {
	n := len(p.Activities)
	start := time.Now()
	code := fmt.Sprintf("filter(Activities, %s)", q)
	log.Debug().
		Str("code", code).
		Msg("filter")
	out, err := expr.Eval(code, p)
	if err != nil {
		return nil, err
	}
	res := out.([]interface{})
	acts := make([]*strava.Activity, len(res))
	for i := range res {
		acts[i] = res[i].(*strava.Activity)
	}
	log.Debug().
		Int("activities {pre}", n).
		Int("activities {post}", len(acts)).
		Dur("elapsed", time.Since(start)).
		Msg("filter")
	return &Pass{Activities: acts, Units: p.Units}, nil
}

// GroupByExpr groups activities by a key
// currently only supports a single key for grouping
func (p *Pass) GroupByExpr(q string) (map[string]*Pass, error) {
	start := time.Now()
	code := fmt.Sprintf("map(Activities, %s)", q)
	log.Debug().
		Str("code", code).
		Msg("groupby")
	out, err := expr.Eval(code, p)
	if err != nil {
		return nil, err
	}
	res := out.([]interface{})
	passes := make(map[string]*Pass, len(res))
	for i, k := range res {
		key, err := cast.ToStringE(k)
		if err != nil {
			return nil, err
		}
		if _, ok := passes[key]; !ok {
			passes[key] = &Pass{Units: p.Units}
		}
		passes[key].Activities = append(passes[key].Activities, p.Activities[i])
	}
	log.Debug().
		Int("activities", len(p.Activities)).
		Int("passes", len(passes)).
		Dur("elapsed", time.Since(start)).
		Msg("groupby")
	return passes, nil
}
