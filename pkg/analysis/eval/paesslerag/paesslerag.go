package paesslerag

import (
	"context"

	"github.com/PaesslerAG/gval"
	"github.com/bzimmer/gravl/pkg/analysis/eval"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

type paesslerag struct {
	language gval.Language
}

func New() eval.Evaluator {
	return &paesslerag{language: gval.Full()}
}

func (p *paesslerag) Filter(ctx context.Context, q string, acts []*strava.Activity) ([]*strava.Activity, error) {
	eval, err := p.language.NewEvaluable(q)
	if err != nil {
		return nil, err
	}
	var res []*strava.Activity
	for i := 0; i < len(acts); i++ {
		ok, err := eval.EvalBool(ctx, acts[i])
		if err != nil {
			return nil, err
		}
		if ok {
			res = append(res, acts[i])
		}
	}
	return res, nil
}

func (p *paesslerag) GroupBy(ctx context.Context, q string, acts []*strava.Activity) (map[string][]*strava.Activity, error) {
	eval, err := p.language.NewEvaluable(q)
	if err != nil {
		return nil, err
	}
	res := make(map[string][]*strava.Activity)
	for i := 0; i < len(acts); i++ {
		key, err := eval.EvalString(ctx, acts[i])
		if err != nil {
			return nil, err
		}
		res[key] = append(res[key], acts[i])
	}
	return res, nil
}
