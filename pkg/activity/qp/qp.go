package qp

import (
	"bytes"
	"context"
	"errors"
	"io"
	"path/filepath"
	"strconv"
	"time"

	"github.com/armon/go-metrics"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	api "github.com/bzimmer/activity"
	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/activity"
	"github.com/bzimmer/gravl/pkg/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/pkg/activity/rwgps"
	"github.com/bzimmer/gravl/pkg/activity/strava"
	"github.com/bzimmer/gravl/pkg/activity/zwift"
)

func poller(c *cli.Context, upd api.Uploader) api.Poller {
	return api.NewPoller(upd,
		api.WithInterval(c.Duration("interval")),
		api.WithIterations(c.Int("iterations")))
}

func exporter(c *cli.Context, name string) (api.Exporter, error) {
	if f, ok := pkg.Runtime(c).Exporters[name]; ok {
		return f(c)
	}
	return nil, errors.New("unknown exporter")
}

func uploader(c *cli.Context, name string) (api.Uploader, error) {
	if f, ok := pkg.Runtime(c).Uploaders[name]; ok {
		return f(c)
	}
	return nil, errors.New("unknown uploader")
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
	i := 0
	for res := range x.poller.Poll(ctx, uploadID) {
		if res.Err != nil {
			return res.Err
		}
		x.metrics.IncrCounter([]string{"upload", "poll"}, 1)
		log.Info().Int("iteration", i).Int64("id", int64(res.Upload.Identifier())).Msg("poll")
		if err := x.encoder.Encode(res.Upload); err != nil {
			return err
		}
		i++
	}
	return nil
}

func upload(c *cli.Context) error {
	fs := pkg.Runtime(c).Fs
	enc := pkg.Runtime(c).Encoder
	met := pkg.Runtime(c).Metrics
	upd, err := uploader(c, c.String("to"))
	if err != nil {
		return err
	}
	x := xfer{
		poller:   poller(c, upd),
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
	upd, err := uploader(c, c.String("to"))
	if err != nil {
		return err
	}
	x := xfer{
		poller:  poller(c, upd),
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
	expr, err := exporter(c, c.String("from"))
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
		Flags:  flags(cfg{from: true, io: true}),
		Action: export,
	}
}

func qp(c *cli.Context) error {
	expr, err := exporter(c, c.String("from"))
	if err != nil {
		return err
	}

	upd, err := uploader(c, c.String("to"))
	if err != nil {
		return err
	}

	x := xfer{
		poller:   poller(c, upd),
		uploader: upd,
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
		Flags:     flags(cfg{from: true, to: true, poll: true, io: true}),
		Action:    qp,
	}
}

type cfg struct {
	from bool
	to   bool
	poll bool
	io   bool
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
				Usage: "Sink data provider"})
	}
	if c.io {
		x = append(x,
			[]cli.Flag{
				&cli.BoolFlag{
					Name:    "overwrite",
					Aliases: []string{"o"},
					Value:   false,
					Usage:   "Overwrite the file if it exists; fail otherwise",
				},
				&cli.StringFlag{
					Name:    "output",
					Aliases: []string{"O"},
					Value:   "",
					Usage:   "The filename to use for writing the contents of the export, if not specified the contents are streamed to stdout",
				},
			}...,
		)
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
		Usage:    "Copy activities from a source to a sink",
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
