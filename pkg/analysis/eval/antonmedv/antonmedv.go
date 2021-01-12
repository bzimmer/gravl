package antonmedv

import (
	"context"
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/spf13/cast"

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

func apply(q string, acts []*strava.Activity) ([]*strava.Activity, error) {
	out, err := expr.Eval(q, &env{Activities: acts})
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

type evaluator struct{}

type env struct {
	Activities []*strava.Activity
}

func New() eval.Evaluator {
	return &evaluator{}
}

func (x *evaluator) Filter(ctx context.Context, q string, acts []*strava.Activity) ([]*strava.Activity, error) {
	code := fmt.Sprintf("filter(Activities, %s)", closure(q))
	return apply(code, acts)
}

func (x *evaluator) Group(ctx context.Context, q string, acts []*strava.Activity) (map[string][]*strava.Activity, error) {
	// map over the activities to generate a group key
	code := fmt.Sprintf("map(Activities, %s)", closure(q))
	out, err := expr.Eval(code, &env{Activities: acts})
	if err != nil {
		return nil, err
	}
	// group all activities into a Group based on their group key
	res := out.([]interface{})
	groups := make(map[string][]*strava.Activity, len(res))
	for i, k := range res {
		var key string
		key, err = cast.ToStringE(k)
		if err != nil {
			return nil, err
		}
		if _, ok := groups[key]; !ok {
			groups[key] = make([]*strava.Activity, 0)
		}
		groups[key] = append(groups[key], acts[i])
	}
	return groups, nil
}
