package qp_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/internal"
	"github.com/bzimmer/gravl/pkg/internal/blackhole"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestWrite(t *testing.T) {
	a := assert.New(t)

	tests := []*internal.Harness{
		{
			Name: "write to app stdout",
			Args: []string{"gravl", "qp", "export", "--from", "blackhole", "61292794933"},
			Before: func(c *cli.Context) error {
				c.App.Writer = bytes.NewBufferString("")
				pkg.Runtime(c).Exporters[blackhole.Provider] = blackhole.ExporterFunc
				return nil
			},
			After: func(c *cli.Context) error {
				bs := c.App.Writer.(*bytes.Buffer)
				a.Equal(blackhole.Data, bs.String())
				return nil
			},
		},
		{
			Name: "write to file",
			Args: []string{"gravl", "qp", "export", "--from", "blackhole", "-O", "/tmp/Foo.gpx", "776765443"},
			Before: func(c *cli.Context) error {
				c.App.Writer = bytes.NewBufferString("")
				pkg.Runtime(c).Exporters[blackhole.Provider] = blackhole.ExporterFunc
				return nil
			},
			After: func(c *cli.Context) error {
				bs := c.App.Writer.(*bytes.Buffer)
				a.Equal(bs.String(), "")
				stat, err := pkg.Runtime(c).Fs.Stat("/tmp/Foo.gpx")
				a.NoError(err)
				a.NotNil(stat)
				a.Equal(int64(len(blackhole.Data)), stat.Size())
				return nil
			},
		},
		{
			Name: "file exists error",
			Args: []string{"gravl", "qp", "export", "--from", "blackhole", "-O", "/tmp/Bar.gpx", "776765443"},
			Before: func(c *cli.Context) error {
				fp, err := pkg.Runtime(c).Fs.Create("/tmp/Bar.gpx")
				a.NoError(err)
				a.NoError(fp.Close())
				pkg.Runtime(c).Exporters[blackhole.Provider] = blackhole.ExporterFunc
				return nil
			},
			Err: os.ErrExist.Error(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, nil, command)
		})
	}
}

func TestList(t *testing.T) {
	a := assert.New(t)

	tests := []*internal.Harness{
		{
			Name: "directory does not exist",
			Args: []string{"gravl", "qp", "list", "/bad/path"},
			Err:  "does not exist",
		},
		{
			Name: "no files",
			Args: []string{"gravl", "qp", "list", "."},
		},
		{
			Name: "one file",
			Args: []string{"gravl", "qp", "list", "/foo/"},
			Before: func(c *cli.Context) error {
				fs := pkg.Runtime(c).Fs
				a.NoError(fs.MkdirAll("/foo/bar/Zwift/Activities", 0777))
				fp, err := fs.Create("/foo/bar/Zwift/Activities/2021-10-01-08:12:13.fit")
				a.NoError(err)
				return fp.Close()
			},
			Counters: map[string]int{
				"gravl.walk.file.attempt":     1,
				"gravl.walk.file.success.fit": 1,
			},
		},
		{
			Name: "two files",
			Args: []string{"gravl", "qp", "list", "/foo/"},
			Before: func(c *cli.Context) error {
				fs := pkg.Runtime(c).Fs
				a.NoError(fs.MkdirAll("/foo/bar/Zwift/Activities", 0777))
				for _, fn := range []string{
					"/foo/bar/Zwift/Activities/2021-10-01-08:12:13.fit",
					"/foo/bar/baz/NotAnActivity.txt",
				} {
					fp, err := fs.Create(fn)
					a.NoError(err)
					err = fp.Close()
					a.NoError(err)
				}
				return nil
			},
			Counters: map[string]int{
				"gravl.walk.file.attempt":              2,
				"gravl.walk.file.success.fit":          1,
				"gravl.walk.file.skipping.unsupported": 1,
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
