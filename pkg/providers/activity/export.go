package activity

//go:generate stringer -type=Format -linecomment -output=export_string.go

import (
	"fmt"
	"io"
	"strings"
)

// Export the contents and metadata about an activity file
type Export struct {
	io.Reader `json:"-"`
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Format    Format `json:"format"`
}

// Format of the file used in exporting and uploading
type Format int

const (
	// Original format (essentially a wildcard)
	Original Format = iota // original
	// GPX format
	GPX // gpx
	// TCX format
	TCX // tcx
	// FIT format
	FIT // fit
)

// MarshalJSON converts a Format enum to a string representation
func (f *Format) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, f.String())), nil
}

// ToFormat converts a file extension (with or without the ".") to a Format
// If no predefined extension exists the Format Original is returned
func ToFormat(format string) Format {
	format = strings.ToLower(format)
	switch format {
	case ".gpx", "gpx":
		return GPX
	case ".tcx", "tcx":
		return TCX
	case ".fit", "fit":
		return FIT
	default:
		return Original
	}
}
