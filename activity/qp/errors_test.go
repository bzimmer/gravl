package qp_test

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	api "github.com/bzimmer/activity"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/internal"
)

var errFlaky = errors.New("flaky provider")

type done struct{}

func (done) Identifier() api.UploadID { return 1 }
func (done) Done() bool               { return true }

// flaky is an uploader and exporter failing whichever operations are set.
type flaky struct {
	upload, status, export, noReader, badReader bool
}

func (f flaky) Upload(context.Context, *api.File) (api.Upload, error) {
	if f.upload {
		return nil, errFlaky
	}
	return done{}, nil
}

func (f flaky) Status(context.Context, api.UploadID) (api.Upload, error) {
	if f.status {
		return nil, errFlaky
	}
	return done{}, nil
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, errFlaky }

func (f flaky) Export(_ context.Context, id int64) (*api.Export, error) {
	if f.export {
		return nil, errFlaky
	}
	var r io.Reader = strings.NewReader("<gpx/>")
	switch {
	case f.noReader:
		r = nil
	case f.badReader:
		r = failingReader{}
	}
	return &api.Export{ID: id, File: &api.File{Name: "x.gpx", Reader: r, Format: api.FormatGPX}}, nil
}

// with registers f as the "flaky" provider, optionally after setup.
func with(f flaky, setup ...cli.BeforeFunc) cli.BeforeFunc {
	return func(c *cli.Context) error {
		rt := gravl.Runtime(c)
		rt.Uploaders["flaky"] = func(*cli.Context) (api.Uploader, error) { return f, nil }
		rt.Exporters["flaky"] = func(*cli.Context) (api.Exporter, error) { return f, nil }
		for _, s := range setup {
			if err := s(c); err != nil {
				return err
			}
		}
		return nil
	}
}

// files creates FIT files for upload.
func files(names ...string) cli.BeforeFunc {
	return func(c *cli.Context) error {
		fs := gravl.Runtime(c).Fs
		for _, name := range names {
			if err := afero.WriteFile(fs, name, []byte("fit"), 0o644); err != nil {
				return err
			}
		}
		return nil
	}
}

// openFailingFs walks directories but cannot open the files in them.
type openFailingFs struct{ afero.Fs }

func (f openFailingFs) Open(name string) (afero.File, error) {
	if strings.HasSuffix(name, ".fit") {
		return nil, errFlaky
	}
	return f.Fs.Open(name)
}

func unopenable(c *cli.Context) error {
	gravl.Runtime(c).Fs = openFailingFs{gravl.Runtime(c).Fs}
	return nil
}

func readOnly(c *cli.Context) error {
	gravl.Runtime(c).Fs = afero.NewReadOnlyFs(gravl.Runtime(c).Fs)
	return nil
}

func TestUploadFailures(t *testing.T) {
	tests := []*internal.Harness{
		{
			Name:   "upload fails",
			Args:   []string{"gravl", "qp", "upload", "--to", "flaky", "/foo"},
			Before: with(flaky{upload: true}, files("/foo/a.fit")),
			Err:    errFlaky.Error(),
		},
		{
			// the failure cancels the walk while it waits to hand over the next file
			Name:   "upload fails with more files queued",
			Args:   []string{"gravl", "qp", "upload", "--to", "flaky", "/foo"},
			Before: with(flaky{upload: true}, files("/foo/a.fit", "/foo/b.fit", "/foo/c.fit")),
			Err:    errFlaky.Error(),
		},
		{
			Name:   "open fails",
			Args:   []string{"gravl", "qp", "upload", "--to", "flaky", "/foo"},
			Before: with(flaky{}, files("/foo/a.fit"), unopenable),
			Err:    errFlaky.Error(),
		},
		{
			Name:   "poll fails",
			Args:   []string{"gravl", "qp", "upload", "--to", "flaky", "--poll", "/foo"},
			Before: with(flaky{status: true}, files("/foo/a.fit")),
			Err:    errFlaky.Error(),
		},
		{
			Name:     "poll succeeds",
			Args:     []string{"gravl", "qp", "upload", "--to", "flaky", "--poll", "/foo"},
			Before:   with(flaky{}, files("/foo/a.fit")),
			Counters: map[string]int{"gravl.upload.poll": 1},
		},
		{
			Name:   "encoding fails",
			Args:   []string{"gravl", "qp", "upload", "--to", "flaky", "/foo"},
			Before: with(flaky{}, files("/foo/a.fit"), internal.FailEncoding),
			Err:    internal.ErrEncode.Error(),
		},
		{
			Name:   "encoding a poll fails",
			Args:   []string{"gravl", "qp", "upload", "--to", "flaky", "--poll", "/foo"},
			Before: with(flaky{}, files("/foo/a.fit"), internal.FailEncoding),
			Err:    internal.ErrEncode.Error(),
		},
		{
			Name: "unknown uploader",
			Args: []string{"gravl", "qp", "upload", "--to", "nowhere", "/foo"},
			Err:  "unknown uploader",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, nil, command)
		})
	}
}

func TestStatusFailures(t *testing.T) {
	tests := []*internal.Harness{
		{
			Name:   "invalid id",
			Args:   []string{"gravl", "qp", "status", "--to", "flaky", "not-a-number"},
			Before: with(flaky{}),
			Err:    "invalid syntax",
		},
		{
			Name:   "status fails",
			Args:   []string{"gravl", "qp", "status", "--to", "flaky", "88191"},
			Before: with(flaky{status: true}),
			Err:    errFlaky.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, nil, command)
		})
	}
}

func TestExportFailures(t *testing.T) {
	a := assert.New(t)
	tests := []*internal.Harness{
		{
			Name:   "invalid id",
			Args:   []string{"gravl", "qp", "export", "--from", "flaky", "not-a-number"},
			Before: with(flaky{}),
			Err:    "invalid syntax",
		},
		{
			Name:   "export fails",
			Args:   []string{"gravl", "qp", "export", "--from", "flaky", "1"},
			Before: with(flaky{export: true}),
			Err:    errFlaky.Error(),
		},
		{
			Name:   "nothing to write",
			Args:   []string{"gravl", "qp", "export", "--from", "flaky", "-O", "/tmp/x.gpx", "1"},
			Before: with(flaky{noReader: true}),
			After: func(c *cli.Context) error {
				exists, err := afero.Exists(gravl.Runtime(c).Fs, "/tmp/x.gpx")
				a.NoError(err)
				a.False(exists)
				return nil
			},
		},
		{
			Name:   "create fails",
			Args:   []string{"gravl", "qp", "export", "--from", "flaky", "-O", "/tmp/x.gpx", "1"},
			Before: with(flaky{}, readOnly),
			AnyErr: true,
		},
		{
			Name:   "read fails",
			Args:   []string{"gravl", "qp", "export", "--from", "flaky", "-O", "/tmp/x.gpx", "1"},
			Before: with(flaky{badReader: true}),
			Err:    errFlaky.Error(),
		},
		{
			Name:   "encoding fails",
			Args:   []string{"gravl", "qp", "export", "--from", "flaky", "-O", "/tmp/x.gpx", "1"},
			Before: with(flaky{}, internal.FailEncoding),
			Err:    internal.ErrEncode.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, nil, command)
		})
	}
}

func TestCopyFailures(t *testing.T) {
	tests := []*internal.Harness{
		{
			Name: "unknown exporter",
			Args: []string{"gravl", "qp", "copy", "--from", "nowhere", "--to", "flaky", "1"},
			Err:  "unknown exporter",
		},
		{
			Name:   "unknown uploader",
			Args:   []string{"gravl", "qp", "copy", "--from", "flaky", "--to", "nowhere", "1"},
			Before: with(flaky{}),
			Err:    "unknown uploader",
		},
		{
			Name:   "invalid id",
			Args:   []string{"gravl", "qp", "copy", "--from", "flaky", "--to", "flaky", "not-a-number"},
			Before: with(flaky{}),
			Err:    "invalid syntax",
		},
		{
			Name:   "export fails",
			Args:   []string{"gravl", "qp", "copy", "--from", "flaky", "--to", "flaky", "1"},
			Before: with(flaky{export: true}),
			Err:    errFlaky.Error(),
		},
		{
			Name:   "upload fails",
			Args:   []string{"gravl", "qp", "copy", "--from", "flaky", "--to", "flaky", "1"},
			Before: with(flaky{upload: true}),
			Err:    errFlaky.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, nil, command)
		})
	}
}

// statFailingFs fails every Stat after the root, as an unreadable subdirectory would.
type statFailingFs struct{ afero.Fs }

func (f statFailingFs) Stat(name string) (os.FileInfo, error) {
	if name == "/foo" {
		return f.Fs.Stat(name)
	}
	return nil, errFlaky
}

func TestListFailures(t *testing.T) {
	tt := &internal.Harness{
		Name: "walk fails partway",
		Args: []string{"gravl", "qp", "list", "/foo"},
		Before: func(c *cli.Context) error {
			if err := files("/foo/a.fit")(c); err != nil {
				return err
			}
			gravl.Runtime(c).Fs = statFailingFs{gravl.Runtime(c).Fs}
			return nil
		},
		Err: errFlaky.Error(),
	}
	internal.Run(t, tt, nil, command)
}
