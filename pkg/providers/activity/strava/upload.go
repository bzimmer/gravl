package strava

import (
	"context"

	"github.com/bzimmer/gravl/pkg/providers/activity"
	"github.com/rs/zerolog/log"
)

type uploader struct {
	s *ActivityService
}

func NewUploader(svc *ActivityService) activity.Uploader {
	return &uploader{s: svc}
}

func (u *uploader) Upload(ctx context.Context, file *activity.File) context.Context {
	go func() {
		up, err := u.s.Upload(ctx, file)
		if err != nil {
			log.Error().Err(ctx.Err()).Msg("upload failed")
			return
		}
		c := u.s.Poll(ctx, up.ID)
		for {
			select {
			case <-ctx.Done():
				log.Error().Err(ctx.Err()).Msg("ctx is done")
				return
			case r, ok := <-c:
				switch {
				case !ok:
					return
				case r.Err != nil:
					log.Error().Err(r.Err).Msg("done")
					return
				case r.Upload != nil:
					log.Info().Int64("uploadID", r.Upload.ID).Msg("not done")
				}
			}
		}
	}()
	return nil
}
