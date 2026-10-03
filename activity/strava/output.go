package strava

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl"
)

// outputFlags returns the flags for writing each result of a per-id command to its own file.
// The names follow the export commands' convention: --output/-O for the destination and
// --overwrite/-o to replace an existing file.
func outputFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"O"},
			Usage: "Write each result to its own file, the path a Go template over the result " +
				`(eg, 'archive/{{.StartDateLocal.Format "2006-01"}}/{{.ID}}.json'); ` +
				"an existing file is skipped unless --overwrite is set",
		},
		&cli.BoolFlag{
			Name:    "overwrite",
			Aliases: []string{"o"},
			Value:   false,
			Usage:   "With --output, replace an existing file instead of skipping it",
		},
	}
}

// fileEncoder writes each encoded value to the file named by rendering a template against it.
// Every file is written to a temporary sibling and renamed into place, so a reader never sees
// a partial file.
type fileEncoder struct {
	fs        afero.Fs
	path      *template.Template
	overwrite bool
	onSkip    func(path string)
}

func newFileEncoder(c *cli.Context) (*fileEncoder, error) {
	tmpl, err := template.New("output").Option("missingkey=error").Parse(c.String("output"))
	if err != nil {
		return nil, fmt.Errorf("parsing --output template: %w", err)
	}
	return &fileEncoder{
		fs:        gravl.Runtime(c).Fs,
		path:      tmpl,
		overwrite: c.Bool("overwrite"),
		onSkip: func(path string) {
			gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, c.Command.Name, "skipped"}, 1)
			log.Info().Str("path", path).Msg("exists, skipping")
		},
	}, nil
}

func (f *fileEncoder) Encode(v any) error {
	var buf bytes.Buffer
	if err := f.path.Execute(&buf, v); err != nil {
		return fmt.Errorf("rendering --output template: %w", err)
	}
	path := buf.String()
	if path == "" {
		return errors.New("--output template rendered an empty path")
	}
	if !f.overwrite {
		switch _, err := f.fs.Stat(path); {
		case err == nil:
			if f.onSkip != nil {
				f.onSkip(path)
			}
			return nil
		case !errors.Is(err, os.ErrNotExist):
			return err
		}
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return writeAtomic(f.fs, path, append(data, '\n'))
}

// writeAtomic writes data to a temporary file beside path, then renames it into place.
func writeAtomic(fs afero.Fs, path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := fs.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := afero.TempFile(fs, dir, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	_, err = tmp.Write(data)
	if cerr := tmp.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = fs.Rename(name, path)
	}
	if err != nil {
		_ = fs.Remove(name)
		return err
	}
	log.Info().Str("path", path).Msg("wrote")
	return nil
}

// readIDs returns one id per non-blank line of r.
func readIDs(r io.Reader) ([]string, error) {
	var ids []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			if _, err := strconv.ParseInt(line, 0, 64); err != nil {
				return nil, err
			}
			ids = append(ids, line)
		}
	}
	return ids, scanner.Err()
}
