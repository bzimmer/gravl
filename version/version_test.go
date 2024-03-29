package version_test

import (
	"testing"

	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/internal"
	"github.com/bzimmer/gravl/version"
)

func command(_ *testing.T, _ string) *cli.Command {
	return version.Command()
}

func TestVersion(t *testing.T) {
	tests := []*internal.Harness{
		{
			Name: "version",
			Args: []string{"gravl", "version"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, nil, command)
		})
	}
}
