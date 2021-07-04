package activity

import (
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/activity"
)

// Write the contents of the export to a file if `output` is specified, else `stdout`
// If the file is written to `output` then the metadata is written to `stdout`, else
// only the file is written to `stdout`.
func Write(c *cli.Context, exp *activity.Export) error {
	if exp == nil || exp.Reader == nil {
		return nil
	}
	// if neither overwrite or output is set use stdout
	if !c.IsSet("overwrite") && !c.IsSet("output") {
		_, err := io.Copy(c.App.Writer, exp)
		return err
	}
	var err error
	var fp *os.File
	var filename = exp.Name
	if c.IsSet("output") {
		// if output is set then use the filename provided by the activity source
		filename = c.String("output")
	}
	// if the file exists and overwrite is not set then error
	if _, err = os.Stat(filename); err == nil && !c.Bool("overwrite") {
		log.Error().Str("filename", filename).Msg("file exists and -o flag not specified")
		return os.ErrExist
	}
	fp, err = os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err = fp.Close(); err != nil {
			log.Error().Err(err).Msg("write")
		}
	}()
	_, err = io.Copy(fp, exp)
	if err != nil {
		return err
	}
	return encoding.Encode(exp)
}

// CollectFunc returns true if the file should be uploaded, false otherwise
type CollectFunc func(path string, info os.FileInfo) bool

// Collect data files ready for uploading to an activity service
// Only files of the format FIT, GPX, or TCX will be considered for uploading
func Collect(name string, f CollectFunc) ([]*activity.File, error) {
	var files []*activity.File
	err := filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		format := activity.ToFormat(filepath.Ext(path))
		switch format {
		case activity.FIT, activity.GPX, activity.TCX:
			// no processing necessary
		case activity.Original:
			log.Info().Str("file", path).Msg("skipping")
			return nil
		}
		if f != nil {
			if !f(path, info) {
				return nil
			}
		}
		log.Info().Str("file", path).Msg("collecting")
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		files = append(files, &activity.File{FQPN: path, Name: filepath.Base(path), Format: format, Reader: file})
		return nil
	})
	return files, err
}
