package analysis

import (
	"flag"
	"testing"

	"github.com/bzimmer/gravl/pkg/analysis/passes/cluster"
	"github.com/bzimmer/gravl/pkg/analysis/passes/totals"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestAnalyzers(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	ssf := &cli.StringSliceFlag{
		Name:       "analyzer",
		HasBeenSet: true,
		Value:      cli.NewStringSlice("totals", "cluster,clusters=5"),
	}
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	a.NoError(ssf.Apply(fs))

	Add(nil, false)
	Add(totals.New(), true)
	Add(cluster.New(), true)

	app := cli.NewApp()
	c := cli.NewContext(app, fs, nil)
	y, err := analyzers(c)
	a.NoError(err)
	a.NotNil(y)
}
