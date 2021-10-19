package qp_test

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bzimmer/activity"
	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/activity/qp"
	"github.com/bzimmer/gravl/pkg/internal"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestOutputOverwrite(t *testing.T) {
	t.Skipf("replace with unittest")
	tests := []struct {
		name                   string
		err, output, overwrite bool
	}{
		{name: "no-args"},
		{name: "only-output", output: true},
		{name: "overwrite-and-output", output: true, overwrite: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			writer := &bytes.Buffer{}
			app := &cli.App{
				Writer: writer,
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "overwrite", Value: false},
					&cli.StringFlag{Name: "output", Value: ""},
				},
				Action: func(c *cli.Context) error {
					exp := &activity.Export{File: &activity.File{Name: tt.name, Reader: strings.NewReader(tt.name)}}
					return qp.Write(c, exp)
				},
				Metadata: map[string]interface{}{
					"enc": pkg.JSON(ioutil.Discard, true),
				},
			}
			var args = []string{""}
			if tt.overwrite {
				args = append(args, "--overwrite")
			}
			if tt.output {
				dirname, err := ioutil.TempDir("", "TestOutputOverwrite")
				a.NoError(err)
				token, err := pkg.Token(16)
				a.NoError(err)
				args = append(args, "--output", filepath.Join(dirname, token))
			}
			a.NoError(app.Run(args))
			if tt.overwrite || tt.output {
				a.Equal("", writer.String())
			} else {
				a.Equal(tt.name, writer.String())
			}
			if tt.output && !tt.overwrite {
				a.Error(app.Run(args))
			}
		})
	}
}

func TestList(t *testing.T) {
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
				a := assert.New(t)
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
				a := assert.New(t)
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
