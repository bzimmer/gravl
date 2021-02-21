package internal_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bzimmer/gravl/pkg/commands"
	"github.com/bzimmer/gravl/pkg/commands/activity/internal"
	"github.com/bzimmer/gravl/pkg/providers/activity"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestWrite(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
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
			writer := &bytes.Buffer{}
			app := &cli.App{
				Writer: writer,
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "overwrite", Value: false},
					&cli.StringFlag{Name: "output", Value: ""},
				},
				Action: func(c *cli.Context) error {
					exp := &activity.Export{Name: tt.name, Reader: strings.NewReader(tt.name)}
					return internal.Write(c, exp)
				},
			}
			var args = []string{""}
			if tt.overwrite {
				args = append(args, "--overwrite")
			}
			if tt.output {
				dirname, err := ioutil.TempDir("", "TestWrite")
				a.NoError(err)
				token, err := commands.Token(16)
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

func TestCollect(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	cwd, err := os.Getwd()
	a.NoError(err)
	files, err := internal.Collect(cwd, nil)
	a.NoError(err)
	a.Equal(0, len(files))

	files, err = internal.Collect(fmt.Sprintf("/tmp/does/not/exist/now/%s", time.Now()), nil)
	a.Error(err)
	a.Nil(files)

	dirname, err := ioutil.TempDir("", "TestCollect")
	a.NoError(err)
	token, err := commands.Token(16)
	a.NoError(err)
	for _, ext := range []string{".fit", ".gpx", ".txt", ".tcx"} {
		f, err := os.Create(filepath.Join(dirname, token+ext))
		a.NoError(err)
		a.NoError(f.Close())
	}

	tests := []struct {
		name    string
		formats map[activity.Format]bool
	}{
		{name: "include none", formats: map[activity.Format]bool{}},
		{name: "include FIT", formats: map[activity.Format]bool{activity.FIT: true}},
		{name: "include TCX", formats: map[activity.Format]bool{activity.TCX: true}},
		{name: "include FIT && GPX", formats: map[activity.Format]bool{activity.FIT: true, activity.GPX: true}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			files, err := internal.Collect(dirname, func(path string, info os.FileInfo) bool {
				ext := filepath.Ext(path)
				format := activity.ToFormat(ext)
				a.NotEqual(activity.Original, format)
				return tt.formats[format]
			})
			a.NoError(err)
			for _, f := range files {
				a.NoError(f.Close())
			}
			a.Equal(len(tt.formats), len(files))
			for i := 0; i < len(files); i++ {
				a.True(tt.formats[files[i].Format])
			}
		})
	}
}