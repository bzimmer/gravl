package web

//go:generate stringer -type=Format -linecomment -output=export_string.go

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// Format of the exported file
type Format int

const (
	// Original format uploaded to strava
	Original Format = iota // original
	// GPX format from an uploaded activity to strava
	GPX // gpx
	// TCX format from an uploaded activity to strava
	TCX // tcx
)

// MarshalJSON converts a Format enum to a string representation
func (f *Format) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, f.String())), nil
}

// ExportFile metadata about an exported file
type ExportFile struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Format    Format    `json:"format"`
	Extension string    `json:"ext"`
	Reader    io.Reader `json:"-"`
}

// ExportService is the API for export endpoints
type ExportService service

// ToFormat converts a Format enum to a string for Strava
func ToFormat(format string) Format {
	format = strings.ToLower(format)
	switch format {
	case "gpx":
		return GPX
	case "tcx":
		return TCX
	default:
		return Original
	}
}

// Export requests the data file for the activity
func (s *ExportService) Export(ctx context.Context, activityID int64, format Format) (*ExportFile, error) {
	uri := fmt.Sprintf("activities/%d/export_%s", activityID, format)
	req, err := s.client.newWebRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	log.Info().
		Str("format", format.String()).
		Int64("activityID", activityID).
		Msg("export")
	res, err := s.client.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return nil, err
		}
	}
	defer res.Body.Close()
	out := &bytes.Buffer{}
	_, err = io.Copy(out, res.Body)
	if err != nil {
		return nil, err
	}
	disposition := res.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(disposition)
	if err != nil {
		return nil, err
	}
	name := params["filename"]
	ext := strings.TrimPrefix(filepath.Ext(name), ".")
	return &ExportFile{
		ID:        activityID,
		Name:      params["filename"],
		Reader:    out,
		Format:    format,
		Extension: ext,
	}, nil
}
