package gravl_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/internal"
)

func TestRuntime(t *testing.T) {
	tests := []*internal.Harness{
		{
			Name: "afters",
			Args: []string{"gravl", "afters"},
			After: gravl.Afters(
				func(c *cli.Context) error {
					gravl.Runtime(c).Metrics.IncrCounter([]string{"afters", "test"}, 1)
					gravl.Runtime(c).Metrics.AddSample([]string{"elapsed"}, float32(time.Since(gravl.Runtime(c).Start).Seconds()))
					return nil
				},
			),
			Counters: map[string]int{
				"gravl.afters.test": 1,
			},
		},
		{
			Name: "token",
			Args: []string{"gravl", "token"},
			Action: func(c *cli.Context) error {
				k, err := gravl.Token(16)
				if err != nil {
					return err
				}
				log.Info().Str("token", k).Msg(c.Command.Name)
				gravl.Runtime(c).Metrics.IncrCounter([]string{"token", "test"}, 1)
				return nil
			},
			Counters: map[string]int{
				"gravl.token.test": 1,
			},
		},
		{
			Name: "befores",
			Args: []string{"gravl", "befores"},
			Before: gravl.Befores(
				func(c *cli.Context) error {
					t, err := gravl.Token(16)
					if err != nil {
						return err
					}
					enc := json.NewEncoder(c.App.Writer)
					if err = enc.Encode(t); err != nil {
						return err
					}
					gravl.Runtime(c).Metrics.IncrCounter([]string{"befores", "test"}, 1)
					return nil
				},
			),
			Counters: map[string]int{
				"gravl.befores.test": 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, nil, func(_ *testing.T, _ string) *cli.Command {
				return &cli.Command{Name: tt.Name, Action: tt.Action}
			})
		})
	}
}

func TestBeforesAftersErrors(t *testing.T) {
	a := assert.New(t)
	failure := errors.New("hook failed")
	var calls []string
	hook := func(name string, err error) func(*cli.Context) error {
		return func(*cli.Context) error {
			calls = append(calls, name)
			return err
		}
	}
	tests := []struct {
		name     string
		run      func(*cli.Context) error
		expected []string
	}{
		{
			name:     "befores skip nil and stop at the first error",
			run:      gravl.Befores(nil, hook("one", nil), hook("two", failure), hook("three", nil)),
			expected: []string{"one", "two"},
		},
		{
			name:     "afters skip nil and stop at the first error",
			run:      gravl.Afters(hook("one", nil), nil, hook("two", failure), hook("three", nil)),
			expected: []string{"one", "two"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(_ *testing.T) {
			calls = nil
			a.ErrorIs(tc.run(nil), failure)
			a.Equal(tc.expected, calls)
		})
	}
}

func TestToken(t *testing.T) {
	a := assert.New(t)
	x, err := gravl.Token(16)
	a.NoError(err)
	y, err := gravl.Token(16)
	a.NoError(err)
	a.Len(x, 24)
	a.NotEqual(x, y)
}
