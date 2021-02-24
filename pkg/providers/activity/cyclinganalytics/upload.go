package cyclinganalytics

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/bzimmer/gravl/pkg/providers/activity"
)

type uploader struct {
	s *RidesService
}

func newUploader(s *RidesService) activity.Uploader {
	return &uploader{s: s}
}

func (u *uploader) Upload(ctx context.Context, file *activity.File) <-chan *activity.Upload {
	res := make(chan *activity.Upload)
	go func() {
		defer close(res)
		var up *activity.Upload
		p, err := u.s.Upload(ctx, file)
		switch {
		case err != nil:
			up = &activity.Upload{Err: err}
		default:
			up = &activity.Upload{Upload: p}
		}
		select {
		case <-ctx.Done():
			log.Debug().Err(ctx.Err()).Msg("ctx is done")
			return
		case res <- up:
			if p.Done() {
				return
			}
		}

		c := u.s.Poll(ctx, p.ID)
		for {
			select {
			case <-ctx.Done():
				log.Debug().Err(ctx.Err()).Msg("ctx is done")
				return
			case x, ok := <-c:
				if !ok {
					return
				}
				switch {
				case x.Err != nil:
					up = &activity.Upload{Err: x.Err}
				case x.Upload != nil:
					up = &activity.Upload{Upload: x.Upload}
				}
			}
			select {
			case <-ctx.Done():
				log.Debug().Err(ctx.Err()).Msg("ctx is done")
				return
			case res <- up:
			}
		}
	}()
	return res
}
