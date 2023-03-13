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
	api "github.com/bzimmer/activity"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/activity"
	"github.com/bzimmer/gravl/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/activity/rwgps"
	"github.com/bzimmer/gravl/activity/strava"
	"github.com/bzimmer/gravl/activity/zwift"
)

func poller(c *cli.Context, upd api.Uploader) api.Poller {
	return api.NewPoller(upd,
		api.WithInterval(c.Duration("interval")),
		api.WithIterations(c.Int("iterations")))
}

func exporter(c *cli.Context, name string) (api.Exporter, error) {
	if f, ok := gravl.Runtime(c).Exporters[name]; ok {
		return f(c)
	}
	return nil, errors.New("unknown exporter")
}

func uploader(c *cli.Context, name string) (api.Uploader, error) {
	if f, ok := gravl.Runtime(c).Uploaders[name]; ok {
		return f(c)
	}
	return nil, errors.New("unknown uploader")
}

func providers(c *cli.Context) error {
	type available struct {
		Exporters []string `json:"exporters"`
		Uploaders []string `json:"uploaders"`
	}
	res := &available{
		Exporters: []string{},
		Uploaders: []string{},
	}
	for key := range gravl.Runtime(c).Exporters {
		res.Exporters = append(res.Exporters, key)
	}
	for key := range gravl.Runtime(c).Uploaders {
		res.Uploaders = append(res.Uploaders, key)
	}
	met := gravl.Runtime(c).Metrics
	met.IncrCounter([]string{c.Command.Name, "exporters"}, float32(len(res.Exporters)))
	met.IncrCounter([]string{c.Command.Name, "uploaders"}, float32(len(res.Uploaders)))
	log.Info().Strs("exporters", res.Exporters).Strs("uploaders", res.Uploaders).Msg(c.Command.Name)
	return gravl.Runtime(c).Encoder.Encode(res)
}

func providersCommand() *cli.Command {
	return &cli.Command{
		Name:        "providers",
		Description: "Return the set of active exporters and uploaders",
		Action:      providers,
	}
}

type xfer struct {
	metrics  *metrics.Metrics
	uploader api.Uploader
	poller   api.Poller
	encoder  gravl.Encoder
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
	fs := gravl.Runtime(c).Fs
	enc := gravl.Runtime(c).Encoder
	met := gravl.Runtime(c).Metrics
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
		var fp afero.File
		fp, err = fs.Open(res.path)
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
		var u api.Upload
		u, err = x.upload(ctx, file)
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
			if err = up(res); err != nil {
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
		encoder: gravl.Runtime(c).Encoder,
		metrics: gravl.Runtime(c).Metrics,
	}
	for i := 0; i < args.Len(); i++ {
		err = func() error {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			var uploadID int64
			uploadID, err = strconv.ParseInt(args.Get(i), 0, 64)
			if err != nil {
				return err
			}
			if err = x.poll(ctx, api.UploadID(uploadID)); err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}

func statusCommand() *cli.Command {
	return &cli.Command{
		Name:        "status",
		Usage:       "Check the status of the upload",
		ArgsUsage:   "UPLOAD_ID (...)",
		Flags:       flags(cfg{to: true, poll: true}),
		Description: "Check the status of an upload",
		Action:      status,
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
		Name:        "list",
		ArgsUsage:   "{FILE | DIRECTORY} (...)",
		Description: "List the files suitable for uploading",
		Action:      list,
	}
}

func export(c *cli.Context) error {
	expr, err := exporter(c, c.String("from"))
	if err != nil {
		return err
	}
	met := gravl.Runtime(c).Metrics
	for i := 0; i < c.NArg(); i++ {
		var exp *api.Export
		var activityID int64
		activityID, err = strconv.ParseInt(c.Args().Get(i), 10, 64)
		if err != nil {
			return err
		}
		exp, err = expr.Export(c.Context, activityID)
		if err != nil {
			return err
		}
		met.IncrCounter([]string{"export", "success"}, 1)
		if err = write(c, exp); err != nil {
			return err
		}
	}
	return nil
}

func exportCommand() *cli.Command {
	return &cli.Command{
		Name:        "export",
		Flags:       flags(cfg{from: true, io: true}),
		Description: "Export an activity from the source",
		Action:      export,
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
		encoder:  gravl.Runtime(c).Encoder,
		metrics:  gravl.Runtime(c).Metrics,
	}
	dur := c.Duration("timeout")
	grp, ctx := errgroup.WithContext(c.Context)
	for i := 0; i < c.NArg(); i++ {
		var activityID int64
		activityID, err = strconv.ParseInt(c.Args().Get(i), 10, 64)
		if err != nil {
			return err
		}
		grp.Go(func() error {
			var cancel func()
			ctx, cancel = context.WithTimeout(ctx, dur)
			defer cancel()
			var u api.Upload
			var exp *api.Export
			exp, err = expr.Export(ctx, activityID)
			if err != nil {
				return err
			}
			log.Info().Int64("id", activityID).Str("exp", exp.Name).Msg("export")
			u, err = x.upload(ctx, exp.File)
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
		Name:        "copy",
		ArgsUsage:   "--from <exporter> --to <uploader> id [id, ...]",
		Flags:       flags(cfg{from: true, to: true, poll: true, io: true}),
		Description: "Copy an activity from a source to a destination",
		Action:      qp,
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
		Usage:    "Manage the flow of activity between different platforms",
		Flags: func() []cli.Flag {
			var x []cli.Flag
			for _, q := range [][]cli.Flag{
				cyclinganalytics.AuthFlags(),
				rwgps.AuthFlags(),
				strava.AuthFlags(),
				zwift.AuthFlags(),
			} {
				x = append(x, q...)
			}
			return x
		}(),
		Subcommands: []*cli.Command{
			copyCommand(),
			exportCommand(),
			listCommand(),
			statusCommand(),
			uploadCommand(),
			providersCommand(),
		},
	}
}
