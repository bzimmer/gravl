package blackhole_test

import (
	"testing"

	"github.com/bzimmer/activity"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/internal/blackhole"
)

func TestBlackhole(t *testing.T) {
	a := assert.New(t)
	exporter := blackhole.NewExporter()
	exp, err := exporter.Export(t.Context(), 1234)
	a.NoError(err)
	a.NotNil(exp)

	uploader := blackhole.NewUploader()
	file := &activity.File{}
	upload, err := uploader.Upload(t.Context(), file)
	a.NoError(err)
	a.NotNil(upload)

	upload, err = uploader.Status(t.Context(), upload.Identifier())
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
