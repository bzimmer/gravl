package web

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/bzimmer/gravl/pkg/providers/activity"
)

var _ activity.Exporter = &ExportService{}

// ExportService is the API for export endpoints
type ExportService service

// Export requests the data file for the activity
func (s *ExportService) Export(ctx context.Context, activityID int64) (*activity.Export, error) {
	return s.ExportWithFormat(ctx, activityID, activity.Original)
}

// Export requests the data file for the activity for the specified format
func (s *ExportService) ExportWithFormat(ctx context.Context, activityID int64, format activity.Format) (*activity.Export, error) {
	uri := fmt.Sprintf("activities/%d/export_%s", activityID, format)
	req, err := s.client.newWebRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	log.Info().Str("format", format.String()).Int64("activityID", activityID).Msg("export")
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
	if format == activity.Original {
		ext := filepath.Ext(params["filename"])
		format = activity.ToFormat(ext)
	}
	return &activity.Export{
		ID: activityID,
		File: &activity.File{
			Reader: out,
			Name:   params["filename"],
			Format: format},
	}, nil
}
