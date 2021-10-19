package qp_test

import (
	"testing"

	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/activity/qp"
	"github.com/bzimmer/gravl/pkg/internal"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func command(t *testing.T, baseURL string) *cli.Command {
	return qp.Command()
}

func TestUpload(t *testing.T) {
	tests := []*internal.Harness{
		{
			Name: "directory does not exist",
			Args: []string{"gravl", "qp", "upload", "--to", "blackhole", "/bad/path"},
			Err:  "does not exist",
		},
		{
			Name: "one file",
			Args: []string{"gravl", "qp", "upload", "--to", "blackhole", "/foo/"},
			Before: func(c *cli.Context) error {
				a := assert.New(t)
				fs := pkg.Runtime(c).Fs
				a.NoError(fs.MkdirAll("/foo/bar/Zwift/Activities", 0777))
				fp, err := fs.Create("/foo/bar/Zwift/Activities/2021-10-01-08:12:13.fit")
				a.NoError(err)
				return fp.Close()
			},
			Counters: map[string]int{
				"gravl.walk.file.attempt":   1,
				"gravl.walk.file.success":   1,
				"gravl.upload.file.success": 1,
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

func TestStatus(t *testing.T) {
	tests := []*internal.Harness{
		{
			Name: "status",
			Args: []string{"gravl", "qp", "status", "--to", "blackhole", "88191"},
			Counters: map[string]int{
				"gravl.upload.poll": 1,
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

func TestExport(t *testing.T) {
	tests := []*internal.Harness{
		{
			Name: "status",
			Args: []string{"gravl", "qp", "export", "--from", "blackhole", "88191"},
			Counters: map[string]int{
				"gravl.export.success": 1,
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
