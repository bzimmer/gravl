package strava_test

import (
	"errors"
	"math"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/activity/strava"
	"github.com/bzimmer/gravl/internal"
)

var errFs = errors.New("fs failure")

// failingFs fails the named operation and delegates everything else.
type failingFs struct {
	afero.Fs
	stat, mkdir, open, rename bool
}

func (f *failingFs) Stat(name string) (os.FileInfo, error) {
	if f.stat {
		return nil, errFs
	}
	return f.Fs.Stat(name)
}

func (f *failingFs) MkdirAll(path string, perm os.FileMode) error {
	if f.mkdir {
		return errFs
	}
	return f.Fs.MkdirAll(path, perm)
}

func (f *failingFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	if f.open {
		return nil, errFs
	}
	return f.Fs.OpenFile(name, flag, perm)
}

func (f *failingFs) Rename(oldname, newname string) error {
	if f.rename {
		return errFs
	}
	return f.Fs.Rename(oldname, newname)
}

func TestOutputFailures(t *testing.T) {
	a := assert.New(t)
	template := `/archive/{{.StartDateLocal.Format "2006-01"}}/{{.ID}}.json`
	withFs := func(fs *failingFs) cli.BeforeFunc {
		return func(c *cli.Context) error {
			fs.Fs = gravl.Runtime(c).Fs
			gravl.Runtime(c).Fs = fs
			return nil
		}
	}
	noFiles := func(c *cli.Context) error {
		// a failed write leaves neither the file nor a temporary sibling
		exists, err := afero.DirExists(gravl.Runtime(c).Fs, "/archive/2026-10")
		a.NoError(err)
		if exists {
			entries, rerr := afero.ReadDir(gravl.Runtime(c).Fs, "/archive/2026-10")
			a.NoError(rerr)
			a.Empty(entries)
		}
		return nil
	}
	tests := []*internal.Harness{
		{
			Name: "empty path",
			Args: []string{"gravl", "strava", "activity", "-O", "{{if false}}x{{end}}", "12345"},
			Err:  "rendered an empty path",
		},
		{
			Name:   "stat fails",
			Args:   []string{"gravl", "strava", "activity", "-O", template, "12345"},
			Before: withFs(&failingFs{stat: true}),
			Err:    errFs.Error(),
		},
		{
			Name:   "mkdir fails",
			Args:   []string{"gravl", "strava", "activity", "-O", template, "12345"},
			Before: withFs(&failingFs{mkdir: true}),
			Err:    errFs.Error(),
		},
		{
			Name:   "temporary file fails",
			Args:   []string{"gravl", "strava", "activity", "-O", template, "12345"},
			Before: withFs(&failingFs{open: true}),
			After:  noFiles,
			Err:    errFs.Error(),
		},
		{
			Name:   "rename fails",
			Args:   []string{"gravl", "strava", "activity", "-O", template, "12345"},
			Before: withFs(&failingFs{rename: true}),
			After:  noFiles,
			Err:    errFs.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, activityHandler(a), command)
		})
	}
}

func TestEncodeFileUnmarshalable(t *testing.T) {
	a := assert.New(t)
	fs := afero.NewMemMapFs()
	err := strava.EncodeFile(fs, "/out/x.json", false, math.NaN())
	a.Error(err)
	exists, err := afero.Exists(fs, "/out/x.json")
	a.NoError(err)
	a.False(exists)
}
