package manual_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/bzimmer/gravl/pkg/internal"
	"github.com/bzimmer/gravl/pkg/manual"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestManual(t *testing.T) {
	a := assert.New(t)

	tests := []*internal.Harness{
		{
			Name: "manual",
			Args: []string{"gravl", "manual"},
			Before: func(c *cli.Context) error {
				c.App.Writer = &bytes.Buffer{}
				return nil
			},
			After: func(c *cli.Context) error {
				s := c.App.Writer.(*bytes.Buffer).String()
				a.Greater(len(s), 0)
				a.Contains(s, "manual")
				return nil
			},
		},
		{
			Name: "manual",
			Args: []string{"gravl", "man"},
			Before: func(c *cli.Context) error {
				c.App.Writer = &bytes.Buffer{}
				return nil
			},
			After: func(c *cli.Context) error {
				s := c.App.Writer.(*bytes.Buffer).String()
				a.Greater(len(s), 0)
				a.Contains(s, "manual")
				return nil
			},
		},
		{
			Name: "manual (not hidden)",
			Args: []string{"gravl", "manual"},
			Before: func(c *cli.Context) error {
				c.App.Writer = &bytes.Buffer{}
				a.Equal("manual", c.App.Commands[0].Name)
				c.App.Commands[0].Hidden = false
				return nil
			},
			After: func(c *cli.Context) error {
				s := c.App.Writer.(*bytes.Buffer).String()
				a.Greater(len(s), 0)
				a.Contains(s, "* [manual](#manual)")
				return nil
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, nil, func(*testing.T, string) *cli.Command {
				return manual.Command()
			})
		})
	}
}

func TestCommands(t *testing.T) {
	a := assert.New(t)

	tests := []*internal.Harness{
		{
			Name: "commands",
			Args: []string{"gravl", "-e", "json", "commands"},
			Before: func(c *cli.Context) error {
				c.App.Writer = &bytes.Buffer{}
				return nil
			},
			After: func(c *cli.Context) error {
				s := c.App.Writer.(*bytes.Buffer).String()
				var m []string
				a.NoError(json.Unmarshal([]byte(s), &m))
				a.Greater(len(m), 0)
				a.Contains(m, "commands commands")
				return nil
			},
		},
		{
			Name: "commands relative",
			Args: []string{"gravl", "-e", "json", "commands", "--relative"},
			Before: func(c *cli.Context) error {
				c.App.Writer = &bytes.Buffer{}
				return nil
			},
			After: func(c *cli.Context) error {
				var x bool
				var m []string
				s := c.App.Writer.(*bytes.Buffer).String()
				a.NoError(json.Unmarshal([]byte(s), &m))
				a.Greater(len(m), 0)
				for i := 0; !x && i < len(m); i++ {
					x = strings.Contains(m[i], "manual.test commands")
				}
				a.True(x, "did not find `manual.test commands`")
				return nil
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, nil, func(*testing.T, string) *cli.Command {
				return manual.Commands()
			})
		})
	}
}
