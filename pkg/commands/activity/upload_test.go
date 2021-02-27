package activity_test

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"

	actcmd "github.com/bzimmer/gravl/pkg/commands/activity"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/activity"
)

var _ activity.Uploader = &uploader{}

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
	status, statuscnt int
}

func (u *uploader) Upload(ctx context.Context, file *activity.File) (activity.Upload, error) {
	return &upload{}, nil
}

func (u *uploader) Status(ctx context.Context, id activity.UploadID) (activity.Upload, error) {
	defer func() { u.statuscnt++ }()
	return &upload{done: u.status == u.statuscnt}, nil
}

func TestUpload(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	tests := []struct {
		name      string
		err, poll bool
		duration  time.Duration
		args      []string
	}{
		{name: "TdF (good; dryrun)", err: false,
			args: []string{"upload", "-n", "../geo/gpx/testdata/2017-07-13-TdF-Stage18.gpx"}},
		{name: "TdF (good; wetrun)", err: false,
			args: []string{"upload", "../geo/gpx/testdata/2017-07-13-TdF-Stage18.gpx"}},
		{name: "TdF (good; wetrun, poll)", err: false, duration: time.Second * 10,
			args: []string{"upload", "-P", "10ms", "-p", "../geo/gpx/testdata/2017-07-13-TdF-Stage18.gpx"}},
		{name: "TdF (missing)", err: true,
			args: []string{"upload", "-n", "2017-07-13-TdF-Stage18.gpx"}},
		{name: "status", err: false, duration: time.Second * 10,
			args: []string{"upload", "-P", "10ms", "-p", "-s", "82992872789392"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := &cli.App{Writer: ioutil.Discard}
			set := flag.NewFlagSet("test", 0)
			err := set.Parse(tt.args)
			a.NoError(err)

			encoding.Add(encoding.JSON(os.Stdout, true))

			context := cli.NewContext(app, set, nil)
			command := actcmd.UploadCommand(func(c *cli.Context) (activity.Uploader, error) {
				return &uploader{status: 1}, nil
			})
			command.Flags = append(command.Flags, &cli.DurationFlag{Name: "timeout", Value: tt.duration})
			err = command.Run(context)
			switch tt.err {
			case true:
				a.Error(err)
			case false:
				a.NoError(err)
			}
		})
	}
}