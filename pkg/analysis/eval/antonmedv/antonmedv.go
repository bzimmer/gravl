package antonmedv

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antonmedv/expr"

	"github.com/bzimmer/gravl/pkg/analysis/eval"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

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

type ISOWeek struct {
	Year, Week int
}

func (w ISOWeek) String() string {
	return fmt.Sprintf("[%04d %02d]", w.Year, w.Week)
}

func isoweek(t time.Time) ISOWeek {
	year, week := t.ISOWeek()
	return ISOWeek{year, week}
}

func newEnv(acts []*strava.Activity) map[string]interface{} {
	return map[string]interface{}{
		"Activities": acts,
		"isoweek":    isoweek,
	}
}

type evaluator struct{}

func New() eval.Evaluator {
	return &evaluator{}
}

func (x *evaluator) Filter(ctx context.Context, q string, acts []*strava.Activity) ([]*strava.Activity, error) {
	env := newEnv(acts)
	code := fmt.Sprintf("filter(Activities, %s)", closure(q))
	pgrm, err := expr.Compile(code, expr.Env(env))
	if err != nil {
		return nil, err
	}
	out, err := expr.Run(pgrm, env)
	if err != nil {
		return nil, err
	}
	res := out.([]interface{})
	p := make([]*strava.Activity, len(res))
	for i := range res {
		p[i] = res[i].(*strava.Activity)
	}
	return p, nil
}

func (x *evaluator) GroupBy(ctx context.Context, q string, acts []*strava.Activity) (map[string][]*strava.Activity, error) {
	env := newEnv(acts)
	// map over the activities to generate a group key
	code := fmt.Sprintf("map(Activities, %s)", closure(q))
	pgrm, err := expr.Compile(code, expr.Env(env))
	if err != nil {
		return nil, err
	}
	out, err := expr.Run(pgrm, env)
	if err != nil {
		return nil, err
	}
	// group all activities into a slice with their key
	res := out.([]interface{})
	groups := make(map[string][]*strava.Activity, len(res))
	for i, k := range res {
		var key = fmt.Sprintf("%v", k)
		groups[key] = append(groups[key], acts[i])
	}
	return groups, nil
}
