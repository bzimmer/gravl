package activity

//go:generate stringer -type=Format -linecomment -output=xfer_string.go

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// Export the contents and metadata about an activity file
type Export struct {
	*File
	ID int64 `json:"id"`
}

type Exporter interface {
	Export(ctx context.Context, activityID int64) (*Export, error)
}

type UploadID int64

type Upload interface {
	Identifier() UploadID
	Done() bool
}

type UploadResult struct {
	Upload Upload
	Err    error
}

type Uploader interface {
	Upload(ctx context.Context, file *File) (Upload, error)
	Status(ctx context.Context, id UploadID) (Upload, error)
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

type Poller struct {
	Uploader Uploader
}

// Poll the status of an upload
//
// The operation will continue until either it is completed, the context
//  is canceled, or the maximum number of iterations have been exceeded.
func (p *Poller) Poll(ctx context.Context, uploadID UploadID) <-chan *UploadResult {
	iterations := 5
	duration := 2 * time.Second
	res := make(chan *UploadResult)
	go func() {
		defer close(res)
		if p.Uploader == nil {
			return
		}
		i := 0
		for ; i < iterations; i++ {
			var r *UploadResult
			log.Info().Int64("uploadID", int64(uploadID)).Msg("status")
			upload, err := p.Uploader.Status(ctx, uploadID)
			switch {
			case err != nil:
				r = &UploadResult{Err: err}
			default:
				r = &UploadResult{Upload: upload}
			}
			select {
			case <-ctx.Done():
				return
			case res <- r:
				if r.Upload.Done() {
					return
				}
			}
			// wait for a bit to let the processing continue
			select {
			case <-ctx.Done():
				return
			case <-time.After(duration):
			}
		}
		if i == iterations {
			log.Warn().Int("polls", iterations).Msg("exceeded max iterations")
		}
	}()
	return res
}
