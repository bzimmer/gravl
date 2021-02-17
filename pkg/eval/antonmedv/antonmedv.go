package antonmedv

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/martinlindhe/unit"

	"github.com/bzimmer/gravl/pkg/eval"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

func closure(f string) string {
	if f == "" {
		return ""
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

func fahrenheit(c float64) float64 {
	return unit.FromCelsius(c).Fahrenheit()
}

func env(acts []*strava.Activity) map[string]interface{} {
	return map[string]interface{}{
		"Activities": acts,
		"isoweek":    isoweek,
		"F":          fahrenheit,
	}
}

type evaluator struct {
	program *vm.Program
}

func compile(q string) (*evaluator, error) {
	pgm, err := expr.Compile(q, expr.Env(env(nil)))
	if err != nil {
		return nil, err
	}
	return &evaluator{pgm}, nil
}

func Mapper(q string) (eval.Mapper, error) {
	return compile(fmt.Sprintf("map(Activities, %s)", closure(q)))
}

func Filterer(q string) (eval.Filterer, error) {
	return compile(fmt.Sprintf("filter(Activities, %s)", closure(q)))
}

func Evaluator(q string) (eval.Evaluator, error) {
	return compile(fmt.Sprintf("map(Activities, %s)", closure(q)))
}

func (x *evaluator) run(acts []*strava.Activity) ([]interface{}, error) {
	out, err := expr.Run(x.program, env(acts))
	if err != nil {
		return nil, err
	}
	return out.([]interface{}), nil
}

func (x *evaluator) Filter(ctx context.Context, acts []*strava.Activity) ([]*strava.Activity, error) {
	res, err := x.run(acts)
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
	return x.run(acts)
}

func (x *evaluator) Bool(ctx context.Context, act *strava.Activity) (bool, error) {
	res, err := x.Eval(ctx, act)
	if err != nil {
		return false, err
	}
	switch z := res.(type) {
	case bool:
		return res.(bool), nil
	default:
		return false, fmt.Errorf("expected type `bool` found `%z`", z)
	}
}

func (x *evaluator) Eval(ctx context.Context, act *strava.Activity) (interface{}, error) {
	res, err := x.run([]*strava.Activity{act})
	if err != nil {
		return nil, err
	}
	return res[0], nil
}