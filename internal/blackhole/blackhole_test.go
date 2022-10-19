package blackhole_test

import (
	"context"
	"testing"

	"github.com/bzimmer/activity"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/internal/blackhole"
)

func TestBlackhole(t *testing.T) {
	a := assert.New(t)
	exporter := blackhole.NewExporter()
	exp, err := exporter.Export(context.Background(), 1234)
	a.NoError(err)
	a.NotNil(exp)

	uploader := blackhole.NewUploader()
	file := &activity.File{}
	upload, err := uploader.Upload(context.Background(), file)
	a.NoError(err)
	a.NotNil(upload)

	upload, err = uploader.Status(context.Background(), upload.Identifier())
	a.NoError(err)
	a.NotNil(upload)
	a.True(upload.Done())

	c := &cli.Context{}
	a.NoError(blackhole.Before(c))
	exporter, err = blackhole.ExporterFunc(c)
	a.NoError(err)
	a.NotNil(exporter)
	uploader, err = blackhole.UploaderFunc(c)
	a.NoError(err)
	a.NotNil(uploader)
}
