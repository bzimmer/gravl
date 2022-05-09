package qp

import (
	"io"
	"os"
	"path/filepath"

	"github.com/armon/go-metrics"
	"github.com/bzimmer/activity"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl"
)

// write the contents of the export to a file if `output` is specified, else `stdout`
// If the file is written to `output` then the metadata is written to `stdout`, else
// only the file is written to `stdout`.
func write(c *cli.Context, exp *activity.Export) error {
	if exp == nil || exp.Reader == nil {
		return nil
	}
	// if neither overwrite or output is set use stdout
	if !c.IsSet("overwrite") && !c.IsSet("output") {
		_, err := io.Copy(c.App.Writer, exp)
		return err
	}
	var err error
	var fp afero.File
	var filename = exp.Name
	if c.IsSet("output") {
		// if output is set then use the filename provided by the activity source
		filename = c.String("output")
	}
	// if the file exists and overwrite is not set then error
	fs := gravl.Runtime(c).Fs
	if _, err = fs.Stat(filename); err == nil && !c.Bool("overwrite") {
		log.Error().Str("filename", filename).Msg("file exists and -o flag not specified")
		return os.ErrExist
	}
	fp, err = fs.Create(filename)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = io.Copy(fp, exp)
	if err != nil {
		return err
	}
	return gravl.Runtime(c).Encoder.Encode(exp)
}

type walkResult struct {
	err  error
	path string
}

// walkFunc returns true if the file should be uploaded, false otherwise
type walkFunc func(met *metrics.Metrics, path string, info os.FileInfo) bool

func formatWalkFunc(met *metrics.Metrics, path string, info os.FileInfo) bool {
	format := activity.ToFormat(filepath.Ext(path))
	switch format {
	case activity.FormatFIT, activity.FormatGPX, activity.FormatTCX:
		met.IncrCounter([]string{"walk", "file", "success", format.String()}, 1)
		return true
	case activity.FormatOriginal:
		// please the linter
	}
	met.IncrCounter([]string{"walk", "file", "skipping", "unsupported"}, 1)
	return false
}

// walk identifies data files ready for uploading to an activity service
// Only files of the format FIT, GPX, or TCX will be considered for uploading
func walk(c *cli.Context, name string, funcs ...walkFunc) <-chan *walkResult {
	ctx := c.Context
	met := gravl.Runtime(c).Metrics
	files := make(chan *walkResult)
	go func() {
		defer close(files)
		err := afero.Walk(gravl.Runtime(c).Fs, name, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			met.IncrCounter([]string{"walk", "file", "attempt"}, 1)
			for _, f := range funcs {
				if !f(met, path, info) {
					return nil
				}
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case files <- &walkResult{path: path}:
				met.IncrCounter([]string{"walk", "file", "success"}, 1)
			}
			return nil
		})
		if err != nil {
			select {
			case <-ctx.Done():
				return
			case files <- &walkResult{err: err}:
			}
		}
	}()
	return files
}
