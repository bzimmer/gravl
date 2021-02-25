package activity

import (
	"context"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/activity"
)

type UploaderFunc func(c *cli.Context) (activity.Uploader, error)

func poll(ctx context.Context, uploader activity.Uploader, uploadID activity.UploadID, follow bool) error {
	p := activity.NewPoller(uploader)
	for res := range p.Poll(ctx, uploadID) {
		if res.Err != nil {
			return res.Err
		}
		if err := encoding.Encode(res.Upload); err != nil {
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
	for i := 0; i < args.Len(); i++ {
		files, err := Collect(args.Get(i), nil)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			log.Warn().Msg("no files specified")
			return nil
		}
		for _, file := range files {
			defer file.Close()
			if dryrun {
				log.Info().Str("file", file.Name).Bool("dryrun", dryrun).Msg("uploading")
				if err := encoding.Encode(map[string]interface{}{
					"dryrun": true,
					"file":   file,
				}); err != nil {
					return err
				}
				continue
			}
			log.Info().Str("file", file.Name).Msg("uploading")
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			u, err := uploader.Upload(ctx, file)
			if err != nil {
				return err
			}
			if !c.Bool("poll") {
				return encoding.Encode(u)
			}
			if err := poll(ctx, uploader, u.Identifier(), true); err != nil {
				return err
			}
		}
	}
	return nil
}

func status(c *cli.Context, uploader activity.Uploader) error {
	args := c.Args()
	for i := 0; i < args.Len(); i++ {
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		uploadID, err := strconv.ParseInt(args.Get(i), 0, 64)
		if err != nil {
			return err
		}
		if err := poll(ctx, uploader, activity.UploadID(uploadID), c.Bool("poll")); err != nil {
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
