package analysis_test

import (
	"context"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

type foo struct {
	Double bool
}

func (f *foo) run(ctx *analysis.Context, pass *analysis.Pass) (interface{}, error) {
	n := len(pass.Activities)
	if f.Double {
		n = n * 2
	}
	return n, nil
}

func TestAnalyze(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	p := &analysis.Pass{
		Activities: []*strava.Activity{
			{Distance: 142000, ElevationGain: 30},
			{Distance: 155000, ElevationGain: 23},
			{Distance: 202000, ElevationGain: 85},
		},
	}

	f := &foo{Double: false}
	fs := flag.NewFlagSet("climbing", flag.ExitOnError)
	fs.BoolVar(&f.Double, "double", f.Double, "double the count")

	x := &analysis.Analyzer{
		Name:  "foo",
		Run:   f.run,
		Flags: fs,
	}
	y := analysis.Analysis{
		Args:      []string{},
		Analyzers: []*analysis.Analyzer{x},
	}

	ctx := analysis.WithContext(context.Background(), analysis.Imperial)
	res, err := y.Run(ctx, p)
	a.NoError(err)
	a.NotNil(res)
	u := res[""].(map[string]interface{})
	a.Equal(3, u[x.Name])

	any, err := analysis.NewAnalysis([]*analysis.Analyzer{x}, []string{"foo", "--double"})
	a.NoError(err)
	a.NotNil(res)
	res, err = any.Run(ctx, p)
	a.NoError(err)
	a.NotNil(res)
	u = res[""].(map[string]interface{})
	a.Equal(6, u[x.Name])
}
