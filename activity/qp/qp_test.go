package qp_test

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/activity/qp"
	"github.com/bzimmer/gravl/internal"
	"github.com/bzimmer/gravl/internal/blackhole"
)

func command(_ *testing.T, _ string) *cli.Command {
	return qp.Command()
}

func TestUpload(t *testing.T) {
	tests := []*internal.Harness{
		{
			Name: "directory does not exist",
			Args: []string{"gravl", "qp", "upload", "--to", "blackhole", "/bad/path"},
			Err:  "does not exist",
			Before: func(c *cli.Context) error {
				gravl.Runtime(c).Uploaders[blackhole.Provider] = blackhole.UploaderFunc
				return nil
			},
		},
		{
			Name: "one file",
			Args: []string{"gravl", "qp", "upload", "--to", "blackhole", "/foo/"},
			Before: func(c *cli.Context) error {
				gravl.Runtime(c).Uploaders[blackhole.Provider] = blackhole.UploaderFunc
				a := assert.New(t)
				fs := gravl.Runtime(c).Fs
				a.NoError(fs.MkdirAll("/foo/bar/Zwift/Activities", 0755))
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
			Before: func(c *cli.Context) error {
				gravl.Runtime(c).Uploaders[blackhole.Provider] = blackhole.UploaderFunc
				return nil
			},
			Counters: map[string]int{
				"gravl.upload.poll": 1,
			},
		},
		{
			Name: "unknown uploader",
			Args: []string{"gravl", "qp", "status", "--to", "nowhere", "9988272"},
			Err:  "unknown uploader",
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
			Name: "unknown exporter",
			Args: []string{"gravl", "qp", "export", "--from", "nowhere", "882733"},
			Err:  "unknown exporter",
		},
		{
			Name: "export",
			Args: []string{"gravl", "qp", "export", "--from", "blackhole", "61292794933"},
			Before: func(c *cli.Context) error {
				gravl.Runtime(c).Exporters[blackhole.Provider] = blackhole.ExporterFunc
				return nil
			},
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

func TestCopy(t *testing.T) {
	tests := []*internal.Harness{
		{
			Name: "export",
			Args: []string{"gravl", "qp", "copy", "--from", "blackhole", "--to", "blackhole", "61292794933"},
			Before: func(c *cli.Context) error {
				gravl.Runtime(c).Uploaders[blackhole.Provider] = blackhole.UploaderFunc
				gravl.Runtime(c).Exporters[blackhole.Provider] = blackhole.ExporterFunc
				return nil
			},
			Counters: map[string]int{
				"gravl.upload.file.success": 1,
				"gravl.upload.poll":         1,
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

func TestProviders(t *testing.T) {
	a := assert.New(t)
	tests := []*internal.Harness{
		{
			Name: "providers",
			Args: []string{"gravl", "-j", "qp", "providers"},
			Before: func(c *cli.Context) error {
				c.App.Writer = new(bytes.Buffer)
				gravl.Runtime(c).Exporters[blackhole.Provider] = blackhole.ExporterFunc
				gravl.Runtime(c).Uploaders[blackhole.Provider] = blackhole.UploaderFunc
				return nil
			},
			After: func(c *cli.Context) error {
				reader, ok := c.App.Writer.(io.Reader)
				a.True(ok)
				var res map[string][]string
				dec := json.NewDecoder(reader)
				err := dec.Decode(&res)
				a.NoError(err)
				a.Equal(map[string][]string{"exporters": {"blackhole"}, "uploaders": {"blackhole"}}, res)
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
