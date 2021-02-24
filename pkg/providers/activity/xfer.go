package activity

//go:generate stringer -type=Format -linecomment -output=xfer_string.go

import (
	"context"
	"fmt"
	"io"
	"strings"
)

// Export the contents and metadata about an activity file
type Export struct {
	*File
	ID int64 `json:"id"`
}

type Exporter interface {
	Export(ctx context.Context, activityID int64) (*Export, error)
}

type Upload struct {
	Upload interface{}
	Err    error
}

type Uploader interface {
	Upload(ctx context.Context, file *File) <-chan *Upload
}

// File for uploading
type File struct {
	io.Reader `json:"-"`
	Name      string `json:"name"`
	Format    Format `json:"format"`
}

func (f *File) Close() error {
	if f.Reader == nil {
		return nil
	}
	if x, ok := f.Reader.(io.Closer); ok {
		return x.Close()
	}
	return nil
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
