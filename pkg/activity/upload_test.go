package activity_test

import (
	"context"
	"flag"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/activity"
	actcmd "github.com/bzimmer/gravl/pkg/activity"
	"github.com/bzimmer/gravl/pkg/internal"
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
	t.Skipf("replace with unittest")
	t.Parallel()
	tdf := internal.TdF(t)
	tests := []struct {
		name      string
		err, poll bool
		duration  time.Duration
		args      []string
	}{
		{name: "TdF (good; dryrun)", err: false,
			args: []string{"upload", "-n", tdf}},
		{name: "TdF (good; wetrun)", err: false,
			args: []string{"upload", tdf}},
		{name: "TdF (good; wetrun, poll)", err: false, duration: time.Second * 10,
			args: []string{"upload", "-P", "10ms", "-p", tdf}},
		{name: "TdF (missing)", err: true,
			args: []string{"upload", "-n", "2017-07-13-TdF-Stage18.gpx"}},
		{name: "status", err: false, duration: time.Second * 10,
			args: []string{"upload", "-P", "10ms", "-p", "-s", "82992872789392"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			app := &cli.App{
				Writer:   ioutil.Discard,
				Metadata: map[string]interface{}{},
			}
			set := flag.NewFlagSet("test", 0)
			a.NoError(set.Parse(tt.args))

			ctx := cli.NewContext(app, set, nil)
			command := actcmd.UploadCommand(func(c *cli.Context) (activity.Uploader, error) {
				return &uploader{status: 1}, nil
			})
			command.Flags = append(command.Flags, &cli.DurationFlag{Name: "timeout", Value: tt.duration})
			err := command.Run(ctx)
			switch tt.err {
			case true:
				a.Error(err)
			case false:
				a.NoError(err)
			}
		})
	}
}
