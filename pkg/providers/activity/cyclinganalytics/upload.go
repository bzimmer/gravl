package cyclinganalytics

import (
	"context"

	"github.com/bzimmer/gravl/pkg/providers/activity"
)

// The operation will continue until either it is completed (status != "processing"), the context
//  is canceled, or the maximum number of iterations have been exceeded.
//
// More information can be found at:
//  https://www.cyclinganalytics.com/developer/api#/user/user_id/upload/upload_id

type uploader struct {
	s *RidesService
}

func newUploader(s *RidesService) activity.Uploader {
	return &uploader{s: s}
}

func (u *uploader) Upload(ctx context.Context, file *activity.File) (activity.Upload, error) {
	return u.s.Upload(ctx, file)
}

func (u *uploader) Status(ctx context.Context, id activity.UploadID) (activity.Upload, error) {
	return u.s.Status(ctx, int64(id))
}
