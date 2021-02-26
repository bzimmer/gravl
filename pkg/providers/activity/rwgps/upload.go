package rwgps

import (
	"context"

	"github.com/bzimmer/gravl/pkg/providers/activity"
)

type uploader struct {
	s *TripsService
}

func newUploader(s *TripsService) activity.Uploader {
	return &uploader{s: s}
}

// Upload uploads a file
func (u *uploader) Upload(ctx context.Context, file *activity.File) (activity.Upload, error) {
	return u.s.Upload(ctx, file)
}

// Status returns the processing status of a file
func (u *uploader) Status(ctx context.Context, id activity.UploadID) (activity.Upload, error) {
	return u.s.Status(ctx, int64(id))
}
