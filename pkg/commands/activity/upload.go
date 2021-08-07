package activity

import (
	"context"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/activity"
)

type UploaderFunc func(c *cli.Context) (activity.Uploader, error)

func poller(c *cli.Context, uploader activity.Uploader) activity.Poller {
	return activity.NewPoller(uploader,
		activity.WithInterval(c.Duration("interval")),
		activity.WithIterations(c.Int("iterations")))
}

func poll(ctx context.Context, enc encoding.Encoder, p activity.Poller, uploadID activity.UploadID, follow bool) error {
	for res := range p.Poll(ctx, uploadID) {
		if res.Err != nil {
			return res.Err
		}
		if err := enc.Encode(res.Upload); err != nil {
			return err
		}
		if !follow {
			return nil
		}
	}
	return ctx.Err()
}

func upload(c *cli.Context, uploader activity.Uploader) error {
	args := c.Args()
	dryrun := c.Bool("dryrun")
	enc := encoding.For(c)
	for i := 0; i < args.Len(); i++ {
		files, err := Collect(args.Get(i), nil)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			log.Warn().Msg("no files specified")
		}
		for _, file := range files {
			defer file.Close()
			log.Info().Str("file", file.Name).Bool("dryrun", dryrun).Msg("uploading")
			switch dryrun {
			case true:
				if err := enc.Encode(map[string]interface{}{"dryrun": true, "file": file}); err != nil {
					return err
				}
			case false:
				ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
				defer cancel()
				u, err := uploader.Upload(ctx, file)
				if err != nil {
					return err
				}
				switch c.Bool("poll") {
				case true:
					p := poller(c, uploader)
					if err = poll(ctx, enc, p, u.Identifier(), true); err != nil {
						return err
					}
				case false:
					if err = enc.Encode(u); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func status(c *cli.Context, uploader activity.Uploader) error {
	args := c.Args()
	p := poller(c, uploader)
	enc := encoding.For(c)
	for i := 0; i < args.Len(); i++ {
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		uploadID, err := strconv.ParseInt(args.Get(i), 0, 64)
		if err != nil {
			return err
		}
		if err := poll(ctx, enc, p, activity.UploadID(uploadID), c.Bool("poll")); err != nil {
			return err
		}
	}
	return nil
}

func UploadCommand(f UploaderFunc) *cli.Command {
	return &cli.Command{
		Name:      "upload",
		Aliases:   []string{"u"},
		Usage:     "Upload an activity file",
		ArgsUsage: "{FILE | DIRECTORY} | UPLOAD_ID (...)",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "status",
				Aliases: []string{"s"},
				Value:   false,
				Usage:   "Check the status of the upload",
			},
			&cli.BoolFlag{
				Name:    "poll",
				Aliases: []string{"p"},
				Value:   false,
				Usage:   "Continually check the status of the request until it is completed",
			},
			&cli.BoolFlag{
				Name:    "dryrun",
				Aliases: []string{"n"},
				Value:   false,
				Usage:   "Show the files which would be uploaded but do not upload them",
			},
			&cli.DurationFlag{
				Name:    "interval",
				Aliases: []string{"P"},
				Value:   time.Second * 2,
				Usage:   "The amount of time to wait between polling for an updated status",
			},
			&cli.IntFlag{
				Name:    "iterations",
				Aliases: []string{"N"},
				Value:   5,
				Usage:   "The max number of polling iterations to perform",
			},
		},
		Action: func(c *cli.Context) error {
			uploader, err := f(c)
			if err != nil {
				return err
			}
			if c.Bool("status") {
				return status(c, uploader)
			}
			return upload(c, uploader)
		},
	}
}
