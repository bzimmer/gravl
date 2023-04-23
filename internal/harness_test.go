package internal_test

import (
	"errors"
	"testing"

	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/internal"
)

func command(_ *testing.T, _ string) *cli.Command {
	return &cli.Command{
		Name: "foo",
		Before: func(c *cli.Context) error {
			gravl.Runtime(c).Metrics.IncrCounter([]string{c.Command.Name, "before"}, 1)
			return nil
		},
		After: func(c *cli.Context) error {
			gravl.Runtime(c).Metrics.IncrCounter([]string{c.Command.Name, "after"}, 1)
			return nil
		},
		Action: func(c *cli.Context) error {
			gravl.Runtime(c).Metrics.IncrCounter([]string{c.Command.Name, "action"}, 1)
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
				gravl.Runtime(c).Metrics.IncrCounter([]string{"app", "before", "tt"}, 1)
				return nil
			},
			After: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:     "harness with err",
			Args:     []string{"gravl", "--json", "foo"},
			Err:      "foo err bar",
			Counters: map[string]int{},
			Before: func(c *cli.Context) error {
				gravl.Runtime(c).Metrics.IncrCounter([]string{"app", "before", "tt"}, 1)
				return nil
			},
			After: func(c *cli.Context) error {
				return errors.New("foo err bar")
			},
		},
		{
			Name: "harness no sample value",
			Args: []string{"gravl", "--json", "foo"},
			Err:  "cannot find sample",
			Counters: map[string]int{
				"does.not.exist": 1,
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
