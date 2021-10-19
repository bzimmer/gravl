package qp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	api "github.com/bzimmer/activity"
	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/activity"
	"github.com/bzimmer/gravl/pkg/activity/blackhole"
	"github.com/bzimmer/gravl/pkg/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/pkg/activity/rwgps"
	"github.com/bzimmer/gravl/pkg/activity/strava"
	"github.com/bzimmer/gravl/pkg/activity/zwift"
)

type UploaderFunc func(c *cli.Context) (api.Uploader, error)

func exporter(c *cli.Context) (api.Exporter, error) {
	exp := c.String("from")
	log.Info().Str("provider", exp).Msg("exporter")
	switch strings.ToLower(exp) {
	case "":
		return nil, errors.New("`from` is a required flag")
	case zwift.Provider:
		if err := zwift.Before(c); err != nil {
			return nil, err
		}
		client := pkg.Runtime(c).Zwift
		return client.Exporter(), nil
	case blackhole.Provider:
		if err := blackhole.Before(c); err != nil {
			return nil, err
		}
		return blackhole.NewExporter(), nil
	case strava.Provider:
		return nil, errors.New("strava exporter needs to be rewritten to use gpx file")
	}
	return nil, fmt.Errorf("unknown exporter {%s}", exp)
}

func uploader(c *cli.Context) (api.Uploader, api.Poller, error) {
	var uploader api.Uploader
	upd := c.String("to")
	log.Info().Str("provider", upd).Msg("uploader")
	switch strings.ToLower(upd) {
	case "":
		return nil, nil, errors.New("`to` is a required flag")
	case blackhole.Provider:
		if err := blackhole.Before(c); err != nil {
			return nil, nil, err
		}
		uploader = blackhole.NewUploader()
	case cyclinganalytics.Provider, "ca":
		if err := cyclinganalytics.Before(c); err != nil {
			return nil, nil, err
		}
		uploader = pkg.Runtime(c).CyclingAnalytics.Uploader()
	case rwgps.Provider, "ridewithgps":
		if err := rwgps.Before(c); err != nil {
			return nil, nil, err
		}
		uploader = pkg.Runtime(c).RideWithGPS.Uploader()
	case strava.Provider:
		if err := strava.Before(c); err != nil {
			return nil, nil, err
		}
		uploader = pkg.Runtime(c).Strava.Uploader()
	default:
		return nil, nil, fmt.Errorf("unknown uploader {%s}", upd)
	}
	p := api.NewPoller(uploader,
		api.WithInterval(c.Duration("interval")),
		api.WithIterations(c.Int("iterations")))
	return uploader, p, nil
}

type xfer struct {
	metrics  *metrics.Metrics
	uploader api.Uploader
	poller   api.Poller
	encoder  pkg.Encoder
}

func (x *xfer) upload(ctx context.Context, export *api.File) (api.Upload, error) {
	out := new(bytes.Buffer)
	_, err := io.Copy(out, export)
	if err != nil {
		return nil, err
	}
	file := &api.File{Reader: out, Format: export.Format, Name: export.Name}
	defer file.Close()
	u, err := x.uploader.Upload(ctx, file)
	if err != nil {
		return nil, err
	}
	x.metrics.IncrCounter([]string{"upload", "file", "success"}, 1)
	return u, nil
}

func (x *xfer) poll(ctx context.Context, uploadID api.UploadID) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for res := range x.poller.Poll(ctx, uploadID) {
		if res.Err != nil {
			return res.Err
		}
		x.metrics.IncrCounter([]string{"upload", "poll"}, 1)
		if err := x.encoder.Encode(res.Upload); err != nil {
			return err
		}
	}
	return nil
}

func upload(c *cli.Context) error {
	fs := pkg.Runtime(c).Fs
	enc := pkg.Runtime(c).Encoder
	met := pkg.Runtime(c).Metrics
	upd, plr, err := uploader(c)
	if err != nil {
		return err
	}

	x := xfer{
		poller:   plr,
		metrics:  met,
		encoder:  enc,
		uploader: upd,
	}

	up := func(res *walkResult) error {
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		fp, err := fs.Open(res.path)
		if err != nil {
			return err
		}
		defer fp.Close()
		file := &api.File{
			Name:     filepath.Base(res.path),
			Filename: res.path,
			Reader:   fp,
			Format:   api.ToFormat(filepath.Ext(res.path)),
		}
		u, err := x.upload(ctx, file)
		if err != nil {
			return err
		}
		if c.Bool("poll") {
			return x.poll(ctx, u.Identifier())
		}
		return enc.Encode(u)
	}

	args := c.Args()
	for i := 0; i < args.Len(); i++ {
		results := walk(c, args.Get(i))
		for res := range results {
			if res.err != nil {
				return res.err
			}
			log.Info().Str("file", res.path).Msg("uploading")
			if err := up(res); err != nil {
				return err
			}
		}
	}
	return nil
}

func uploadCommand() *cli.Command {
	return &cli.Command{
		Name:      "upload",
		ArgsUsage: "{FILE | DIRECTORY} (...)",
		Flags:     flags(cfg{to: true, poll: true}),
		Action:    upload,
	}
}

func status(c *cli.Context) error {
	args := c.Args()
	_, plr, err := uploader(c)
	if err != nil {
		return err
	}
	x := xfer{
		poller:  plr,
		encoder: pkg.Runtime(c).Encoder,
		metrics: pkg.Runtime(c).Metrics,
	}
	for i := 0; i < args.Len(); i++ {
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		uploadID, err := strconv.ParseInt(args.Get(i), 0, 64)
		if err != nil {
			return err
		}
		if err := x.poll(ctx, api.UploadID(uploadID)); err != nil {
			return err
		}
	}
	return nil
}

func statusCommand() *cli.Command {
	return &cli.Command{
		Name:      "status",
		Usage:     "Check the status of the upload",
		ArgsUsage: "UPLOAD_ID (...)",
		Flags:     flags(cfg{to: true, poll: true}),
		Action:    status,
	}
}

func list(c *cli.Context) error {
	for i := 0; i < c.NArg(); i++ {
		for path := range walk(c, c.Args().Get(i), formatWalkFunc) {
			if path.err != nil {
				return path.err
			}
			log.Info().Str("file", path.path).Msg(c.Command.Name)
		}
	}
	return nil
}

func listCommand() *cli.Command {
	return &cli.Command{
		Name:      "list",
		ArgsUsage: "{FILE | DIRECTORY} (...)",
		Action:    list,
	}
}

func export(c *cli.Context) error {
	expr, err := exporter(c)
	if err != nil {
		return err
	}
	met := pkg.Runtime(c).Metrics
	for i := 0; i < c.NArg(); i++ {
		activityID, err := strconv.ParseInt(c.Args().Get(i), 10, 64)
		if err != nil {
			return err
		}
		exp, err := expr.Export(c.Context, activityID)
		if err != nil {
			return err
		}
		met.IncrCounter([]string{"export", "success"}, 1)
		if err := Write(c, exp); err != nil {
			return err
		}
	}
	return nil
}

func exportCommand() *cli.Command {
	return &cli.Command{
		Name:   "export",
		Flags:  flags(cfg{from: true}),
		Action: export,
	}
}

func qp(c *cli.Context) error {
	expr, err := exporter(c)
	if err != nil {
		return err
	}

	uplr, plr, err := uploader(c)
	if err != nil {
		return err
	}

	x := xfer{
		poller:   plr,
		uploader: uplr,
		encoder:  pkg.Runtime(c).Encoder,
		metrics:  pkg.Runtime(c).Metrics,
	}
	dur := c.Duration("timeout")
	grp, ctx := errgroup.WithContext(c.Context)
	for i := 0; i < c.NArg(); i++ {
		activityID, err := strconv.ParseInt(c.Args().Get(i), 10, 64)
		if err != nil {
			return err
		}
		grp.Go(func() error {
			var cancel func()
			ctx, cancel = context.WithTimeout(ctx, dur)
			defer cancel()
			exp, err := expr.Export(ctx, activityID)
			if err != nil {
				return err
			}
			log.Info().Int64("id", activityID).Str("exp", exp.Name).Msg("export")
			u, err := x.upload(ctx, exp.File)
			if err != nil {
				return err
			}
			return x.poll(ctx, u.Identifier())
		})
	}
	return grp.Wait()
}

func copyCommand() *cli.Command {
	return &cli.Command{
		Name:      "copy",
		ArgsUsage: "--from <exporter> --to <uploader> id [id, ...]",
		Flags:     flags(cfg{from: true, to: true, poll: true}),
		Action:    qp,
	}
}

type cfg struct {
	from bool
	to   bool
	poll bool
}

func flags(c cfg) []cli.Flag {
	x := activity.RateLimitFlags()
	if c.from {
		x = append(x,
			&cli.StringFlag{
				Name:  "from",
				Usage: "Source data provider"})
	}
	if c.to {
		x = append(x,
			&cli.StringFlag{
				Name:  "to",
				Usage: "Destination data provider"})
	}
	if c.poll {
		x = append(x,
			[]cli.Flag{
				&cli.BoolFlag{
					Name:  "poll",
					Value: false,
					Usage: "Continually check the status of the request until it is completed",
				},
				&cli.DurationFlag{
					Name:  "interval",
					Value: time.Second * 2,
					Usage: "The amount of time to wait between polling for an updated status",
				},
				&cli.IntFlag{
					Name:    "iterations",
					Aliases: []string{"N"},
					Value:   5,
					Usage:   "The max number of polling iterations to perform",
				},
			}...)
	}
	for _, q := range [][]cli.Flag{
		cyclinganalytics.AuthFlags(),
		rwgps.AuthFlags(),
		strava.AuthFlags(),
		zwift.AuthFlags(),
	} {
		x = append(x, q...)
	}
	return x
}

/*
gravl qp list (files, directories)...
gravl qp upload --to <uploader> (files, directories)...

gravl qp export --from <exporter> (ids)...
gravl qp copy --from <exporter> --to <uploader> (ids)...

gravl qp status --from <exporter> (ids)...
*/

func Command() *cli.Command {
	return &cli.Command{
		Name:     "qp",
		Category: "activity",
		Usage:    "Copy activities from a source to a destination",
		Flags:    flags(cfg{}),
		Subcommands: []*cli.Command{
			copyCommand(),
			exportCommand(),
			listCommand(),
			statusCommand(),
			uploadCommand(),
		},
	}
}
