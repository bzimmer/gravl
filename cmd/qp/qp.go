package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	cacmd "github.com/bzimmer/gravl/pkg/commands/activity/cyclinganalytics"
	stcmd "github.com/bzimmer/gravl/pkg/commands/activity/strava"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/commands/gravl"
	"github.com/bzimmer/gravl/pkg/providers/activity"
	"github.com/bzimmer/gravl/pkg/providers/activity/cyclinganalytics"
)

func export(ctx context.Context, c *cli.Context, activityID int64) (*activity.Export, error) {
	client, err := stcmd.NewWebClient(c)
	if err != nil {
		return nil, err
	}
	log.Info().Int64("activityID", activityID).Msg("exporting")
	return client.Export.Export(ctx, activityID, activity.Original)
}

func upload(c *cli.Context, export *activity.Export) error {
	out := new(bytes.Buffer)
	_, err := io.Copy(out, export)
	if err != nil {
		return err
	}
	file := &activity.File{Reader: out, Format: export.Format, Name: export.Name}
	defer file.Close()
	client, err := cacmd.NewClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	upload, err := client.Rides.Upload(ctx, cyclinganalytics.Me, file)
	if err != nil {
		return err
	}
	pc := client.Rides.Poll(ctx, upload.UserID, upload.ID)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case res, ok := <-pc:
			if !ok {
				return nil
			}
			if res.Err != nil {
				return res.Err
			}
			if err := encoding.Encode(res.Upload); err != nil {
				return err
			}
		}
	}
}

func qp(c *cli.Context) error {
	for _, arg := range c.Args().Slice() {
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		actID, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return err
		}
		exp, err := export(ctx, c, actID)
		if err != nil {
			return err
		}
		if err := upload(c, exp); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var flags = gravl.Flags("gravl.yaml")
	flags = append(flags, stcmd.AuthFlags...)
	flags = append(flags, cacmd.AuthFlags...)
	app := &cli.App{
		Name:     "qp",
		HelpName: "qp",
		Usage:    "Copy activities from Strava to CyclingAnalytics",
		Flags:    flags,
		Before:   gravl.Befores(gravl.InitLogging(), gravl.InitEncoding(), gravl.InitConfig()),
		Action:   qp,
		ExitErrHandler: func(c *cli.Context, err error) {
			if err == nil {
				return
			}
			log.Error().Err(err).Msg(c.App.Name)
		},
	}
	if err := app.RunContext(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
