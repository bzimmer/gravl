package activity_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/providers/activity"
)

type upload struct {
	done bool
}

func (u *upload) Identifier() activity.UploadID {
	return activity.UploadID(1122)
}

func (u *upload) Done() bool {
	return u.done
}

type uploader struct {
	err               bool
	status, statuscnt int
}

func (u *uploader) Upload(ctx context.Context, file *activity.File) (activity.Upload, error) {
	return &upload{}, nil
}

func (u *uploader) Status(ctx context.Context, id activity.UploadID) (activity.Upload, error) {
	defer func() { u.statuscnt++ }()
	if u.err {
		return nil, errors.New("uploader error")
	}
	return &upload{done: u.status == u.statuscnt}, nil
}

func TestPoller(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	tests := []struct {
		name       string
		err        bool
		it, status int
		in, to     time.Duration
	}{
		{name: "< iterations", status: 3, it: 5, in: time.Millisecond * 10},
		{name: "max iterations", status: 100, it: 5, in: time.Millisecond * 10},
		{name: "errors", status: 1, it: 5, in: time.Millisecond * 10, err: true},
		{name: "ctx timeout", status: 1, it: 5, in: time.Second, to: time.Millisecond * 10, err: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.to > 0 {
				var cancel func()
				ctx, cancel = context.WithTimeout(ctx, tt.to)
				defer cancel()
			}
			u := &uploader{status: tt.status, err: tt.err}
			p := activity.NewPoller(u, activity.WithInterval(tt.in), activity.WithIterations(tt.it))
			a.NotNil(p)
			for x := range p.Poll(ctx, activity.UploadID(11011)) {
				a.NotNil(x)
				switch tt.err {
				case true:
					a.Nil(x.Upload)
					a.Error(x.Err)
				case false:
					a.NoError(x.Err)
					a.NotNil(x.Upload)
				}
			}
		})
	}
}
