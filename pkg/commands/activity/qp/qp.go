package qp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/commands/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/pkg/commands/activity/rwgps"
	"github.com/bzimmer/gravl/pkg/commands/activity/strava"
	"github.com/bzimmer/gravl/pkg/commands/activity/zwift"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/activity"
)

func poller(c *cli.Context, uploader activity.Uploader) activity.Poller {
	return activity.NewPoller(uploader,
		activity.WithInterval(c.Duration("interval")),
		activity.WithIterations(c.Int("iterations")))
}

func exporter(c *cli.Context) (activity.Exporter, error) {
	exp := c.String("exporter")
	log.Info().Str("provider", exp).Msg("exporter")
	switch exp {
	case "zwift":
		client, err := zwift.NewClient(c)
		if err != nil {
			return nil, err
		}
		return client.NewExporter(), nil
	case "strava": // nolint
		client, err := strava.NewWebClient(c)
		if err != nil {
			return nil, err
		}
		return client.NewExporter(), nil
	}
	return nil, fmt.Errorf("unknown exporter {%s}", exp)
}

func uploader(c *cli.Context) (activity.Uploader, error) {
	upd := c.String("uploader")
	log.Info().Str("provider", upd).Msg("uploader")
	switch upd {
	case "ca", "cyclinganalytics":
		client, err := cyclinganalytics.NewClient(c)
		if err != nil {
			return nil, err
		}
		return client.Uploader(), nil
	case "rwgps":
		client, err := rwgps.NewClient(c)
		if err != nil {
			return nil, err
		}
		return client.Uploader(), nil
	case "strava":
		client, err := strava.NewAPIClient(c)
		if err != nil {
			return nil, err
		}
		return client.Uploader(), nil
	}
	return nil, fmt.Errorf("unknown importer {%s}", upd)
}

func upload(c *cli.Context, upd activity.Uploader, export *activity.Export) error {
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
	u, err := upd.Upload(ctx, file)
	if err != nil {
		return err
	}
	p := poller(c, upd)
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
	exp, err := exporter(c)
	if err != nil {
		return err
	}
	upd, err := uploader(c)
	if err != nil {
		return err
	}
	args := c.Args()
	for i := 0; i < args.Len(); i++ {
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		actID, err := strconv.ParseInt(args.Get(i), 10, 64)
		if err != nil {
			return err
		}
		log.Info().Int64("activityID", actID).Msg("export")
		exp, err := exp.Export(ctx, actID)
		if err != nil {
			return err
		}
		if err := upload(c, upd, exp); err != nil {
			return err
		}
	}
	return nil
}

var flags = func() []cli.Flag {
	f := []cli.Flag{
		&cli.StringFlag{
			Name:    "exporter",
			Aliases: []string{"e"},
			Value:   "",
			Usage:   "Export data provider"},
		&cli.StringFlag{
			Name:    "uploader",
			Aliases: []string{"u"},
			Value:   "",
			Usage:   "Upload data provider"},
	}
	f = append(f, cyclinganalytics.AuthFlags...)
	f = append(f, rwgps.AuthFlags...)
	f = append(f, strava.AuthFlags...)
	f = append(f, zwift.AuthFlags...)
	return f
}()

var Command = &cli.Command{
	Name:      "qp",
	Category:  "activity",
	Usage:     "Copy an activity from an exporter to an uploader",
	ArgsUsage: "ACTIVITY_ID (...)",
	Flags:     flags,
	Action:    qp,
}
