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

// Format of the exported file
type Format int

const (
	// Original format
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

// ToFormat converts a string to a Format enum
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
