package strava

import (
	"context"

	"github.com/bzimmer/gravl/pkg/providers/activity"
)

// The operation will continue until either it is completed, the context
//  is canceled, or the maximum number of iterations have been exceeded.
//
// More information can be found at:
//   https://developers.strava.com/docs/uploads/
//   A successful upload will return a response with an upload ID. You may use this ID to poll the
//   status of your upload. Strava recommends polling no more than once a second. The mean processing
//   time is around 8 seconds.

type uploader struct {
	s *ActivityService
}

func newUploader(s *ActivityService) activity.Uploader {
	return &uploader{s: s}
}

func (u *uploader) Upload(ctx context.Context, file *activity.File) (activity.Upload, error) {
	return u.s.Upload(ctx, file)
}

func (u *uploader) Status(ctx context.Context, id activity.UploadID) (activity.Upload, error) {
	return u.s.Status(ctx, int64(id))
}
