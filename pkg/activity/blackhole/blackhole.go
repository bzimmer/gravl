package blackhole

import (
	"context"
	"strings"

	"github.com/bzimmer/activity"
	"github.com/urfave/cli/v2"
)

const Provider = "blackhole"

type uploadable struct {
	id   activity.UploadID
	done bool
}

func (u *uploadable) Identifier() activity.UploadID {
	return u.id
}

func (u *uploadable) Done() bool {
	return u.done
}

type blackhole struct {
	status, statuscnt int
}

func NewUploader() activity.Uploader {
	return &blackhole{}
}

func NewExporter() activity.Exporter {
	return &blackhole{}
}

// Upload uploads a file
func (b *blackhole) Upload(ctx context.Context, file *activity.File) (activity.Upload, error) {
	return &uploadable{}, nil
}

// Status returns the processing status of a file
func (b *blackhole) Status(ctx context.Context, id activity.UploadID) (activity.Upload, error) {
	defer func() { b.statuscnt++ }()
	return &uploadable{id: id, done: b.status == b.statuscnt}, nil
}

func (b *blackhole) Export(ctx context.Context, activityID int64) (*activity.Export, error) {
	return &activity.Export{
		ID: activityID,
		File: &activity.File{
			Format:   activity.FormatGPX,
			Name:     "Foo",
			Filename: "Foo.gpx",
			Reader:   strings.NewReader("blackhole-foo.gpx"),
		},
	}, nil
}

func Before(c *cli.Context) error {
	return nil
}
