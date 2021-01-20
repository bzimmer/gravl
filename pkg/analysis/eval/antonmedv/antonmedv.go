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

type evaluator struct {
	q string
}

func New(q string) eval.Evaluator {
	return &evaluator{q: closure(q)}
}

func run(q string, acts []*strava.Activity) ([]interface{}, error) {
	env := newEnv(acts)
	pgrm, err := expr.Compile(q, expr.Env(env))
	if err != nil {
		return nil, err
	}
	out, err := expr.Run(pgrm, env)
	if err != nil {
		return nil, err
	}
	return out.([]interface{}), nil
}

func (x *evaluator) Filter(ctx context.Context, acts []*strava.Activity) ([]*strava.Activity, error) {
	code := fmt.Sprintf("filter(Activities, %s)", x.q)
	res, err := run(code, acts)
	if err != nil {
		return nil, err
	}
	p := make([]*strava.Activity, len(res))
	for i := range res {
		p[i] = res[i].(*strava.Activity)
	}
	return p, nil
}

func (x *evaluator) Map(ctx context.Context, acts []*strava.Activity) ([]interface{}, error) {
	code := fmt.Sprintf("map(Activities, %s)", x.q)
	return run(code, acts)
}
