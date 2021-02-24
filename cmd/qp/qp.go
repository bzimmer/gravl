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
)

func exporter(c *cli.Context) (activity.Exporter, error) {
	client, err := stcmd.NewWebClient(c)
	if err != nil {
		return nil, err
	}
	return client.NewExporter(), nil
}

func uploader(c *cli.Context) (activity.Uploader, error) {
	client, err := cacmd.NewClient(c)
	if err != nil {
		return nil, err
	}
	return client.Uploader(), nil
}

func upload(c *cli.Context, uploadr activity.Uploader, export *activity.Export) error {
	out := new(bytes.Buffer)
	_, err := io.Copy(out, export)
	if err != nil {
		return err
	}
	file := &activity.File{Reader: out, Format: export.Format, Name: export.Name}
	defer file.Close()
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	log.Info().Int64("activityID", export.ID).Msg("upload")
	u, err := uploadr.Upload(ctx, file)
	if err != nil {
		return err
	}
	p := activity.Poller{Uploader: uploadr}
	for res := range p.Poll(ctx, u.Identifier()) {
		if res.Err != nil {
			return res.Err
		}
		if err := encoding.Encode(res.Upload); err != nil {
			return err
		}
	}
	return nil
}

func qp(c *cli.Context) error {
	if c.NArg() == 0 {
		help := c.App.Command("help")
		if help == nil {
			return nil
		}
		return help.Run(c)
	}
	expr, err := exporter(c)
	if err != nil {
		return err
	}
	uplr, err := uploader(c)
	if err != nil {
		return err
	}
	for _, arg := range c.Args().Slice() {
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		actID, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return err
		}
		log.Info().Int64("activityID", actID).Msg("export")
		exp, err := expr.Export(ctx, actID)
		if err != nil {
			return err
		}
		if err := upload(c, uplr, exp); err != nil {
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
		Name:      "qp",
		HelpName:  "qp",
		Usage:     "Copy activities from Strava to CyclingAnalytics",
		ArgsUsage: "ACTIVITY_ID (...)",
		Flags:     flags,
		Before:    gravl.Befores(gravl.InitLogging(), gravl.InitEncoding(), gravl.InitConfig()),
		Action:    qp,
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
