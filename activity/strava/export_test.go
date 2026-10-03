package strava

import (
	"text/template"

	"github.com/spf13/afero"
)

// EncodeFile writes v through a fileEncoder, for values no command produces.
func EncodeFile(fs afero.Fs, path string, overwrite bool, v any) error {
	enc := &fileEncoder{fs: fs, path: template.Must(template.New("output").Parse(path)), overwrite: overwrite}
	return enc.Encode(v)
}

// ParseRateLimits exposes the header parsing.
var ParseRateLimits = parseRateLimits //nolint:gochecknoglobals // test export
