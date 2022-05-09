package activity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/activity"
	"github.com/bzimmer/gravl/internal"
)

func TestRateLimitFlags(t *testing.T) {
	a := assert.New(t)
	flags := activity.RateLimitFlags()
	a.Equal(3, len(flags))
}

func TestDateRangeFlags(t *testing.T) {
	a := assert.New(t)
	flags := activity.DateRangeFlags()
	a.Equal(2, len(flags))
}

func TestDateRange(t *testing.T) {
	a := assert.New(t)
	tests := []*internal.Harness{
		{
			Name: "noflags",
			Args: []string{"gravl", "noflags"},
			Action: func(c *cli.Context) error {
				before, after, err := activity.DateRange(c)
				a.NoError(err)
				a.Zero(before)
				a.Zero(after)
				return nil
			},
		},
		{
			Name: "before",
			Args: []string{"gravl", "before", "--before", "yesterday"},

			Action: func(c *cli.Context) error {
				before, after, err := activity.DateRange(c)
				a.NoError(err)
				a.NotZero(before)
				a.Zero(after)
				return nil
			},
		},
		{
			Name: "after",
			Args: []string{"gravl", "after", "--after", "yesterday"},
			Action: func(c *cli.Context) error {
				before, after, err := activity.DateRange(c)
				a.NoError(err)
				a.NotZero(before, "before")
				a.NotZero(after, "after")
				return nil
			},
		},
		{
			Name: "both",
			Args: []string{"gravl", "both", "--after", "two weeks ago", "--before", "yesterday"},
			Action: func(c *cli.Context) error {
				before, after, err := activity.DateRange(c)
				a.NoError(err)
				a.NotZero(before)
				a.NotZero(after)
				return nil
			},
		},
		{
			Name: "err",
			Args: []string{"gravl", "err", "--before", "two weeks ago", "--after", "yesterday"},
			Action: func(c *cli.Context) error {
				before, after, err := activity.DateRange(c)
				a.Error(err)
				a.Zero(before, "before")
				a.Zero(after, "after")
				return nil
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			command := func(t *testing.T, baseURL string) *cli.Command {
				return &cli.Command{
					Name:   tt.Name,
					Flags:  activity.DateRangeFlags(),
					Action: tt.Action,
				}
			}
			internal.Run(t, tt, nil, command)
		})
	}
}
