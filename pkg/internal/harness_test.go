package internal_test

import (
	"testing"

	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/internal"
	"github.com/urfave/cli/v2"
)

func command(t *testing.T, baseURL string) *cli.Command {
	return &cli.Command{
		Name: "foo",
		Before: func(c *cli.Context) error {
			pkg.Runtime(c).Metrics.IncrCounter([]string{c.Command.Name, "before"}, 1)
			return nil
		},
		After: func(c *cli.Context) error {
			pkg.Runtime(c).Metrics.IncrCounter([]string{c.Command.Name, "after"}, 1)
			return nil
		},
		Action: func(c *cli.Context) error {
			pkg.Runtime(c).Metrics.IncrCounter([]string{c.Command.Name, "action"}, 1)
			return nil
		},
	}
}

func TestHarness(t *testing.T) {
	tests := []*internal.Harness{
		{
			Name: "harness",
			Args: []string{"gravl", "foo"},
			Counters: map[string]int{
				"gravl.app.before.tt": 1,
				"gravl.foo.action":    1,
				"gravl.foo.after":     1,
				"gravl.foo.before":    1,
			},
			Before: func(c *cli.Context) error {
				pkg.Runtime(c).Metrics.IncrCounter([]string{"app", "before", "tt"}, 1)
				return nil
			},
			After: func(c *cli.Context) error {
				return nil
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, nil, command)
		})
	}
}
