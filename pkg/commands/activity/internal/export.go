package internal

import (
	"io"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/providers/activity"
)

// Write the contents of the export to a file if `output` is specified, else stdout
func Write(c *cli.Context, exp *activity.Export) error {
	// if neither overwrite or output is set use stdout
	if !c.IsSet("overwrite") && !c.IsSet("output") {
		_, err := io.Copy(os.Stdout, exp.Reader)
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
	defer fp.Close()
	_, err = io.Copy(fp, exp.Reader)
	return err
}
