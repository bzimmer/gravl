package qp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

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
	switch strings.ToLower(exp) {
	case "":
		return nil, errors.New("exporter is a required flag")
	case "zwift":
		client, err := zwift.NewClient(c)
		if err != nil {
			return nil, err
		}
		return client.Exporter(), nil
	case "strava": // nolint
		client, err := strava.NewWebClient(c)
		if err != nil {
			return nil, err
		}
		return client.Exporter(), nil
	}
	return nil, fmt.Errorf("unknown exporter {%s}", exp)
}

func uploader(c *cli.Context) (activity.Uploader, activity.Poller, error) {
	var updr activity.Uploader
	upd := c.String("uploader")
	log.Info().Str("provider", upd).Msg("uploader")
	switch strings.ToLower(upd) {
	case "":
		return nil, nil, errors.New("uploader is a required flag")
	case "ca", "cyclinganalytics":
		client, err := cyclinganalytics.NewClient(c)
		if err != nil {
			return nil, nil, err
		}
		updr = client.Uploader()
	case "rwgps", "ridewithgps":
		client, err := rwgps.NewClient(c)
		if err != nil {
			return nil, nil, err
		}
		updr = client.Uploader()
	case "strava":
		client, err := strava.NewAPIClient(c)
		if err != nil {
			return nil, nil, err
		}
		updr = client.Uploader()
	default:
		return nil, nil, fmt.Errorf("unknown uploader {%s}", upd)
	}
	return updr, poller(c, updr), nil
}

func upload(c *cli.Context, upd activity.Uploader, plr activity.Poller, export *activity.Export) error {
	log.Info().Int64("activityID", export.ID).Msg("upload")
	out := new(bytes.Buffer)
	_, err := io.Copy(out, export)
	if err != nil {
		return err
	}
	file := &activity.File{Reader: out, Format: export.Format, Name: export.Name}
	defer file.Close()
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	u, err := upd.Upload(ctx, file)
	if err != nil {
		return err
	}
	for res := range plr.Poll(ctx, u.Identifier()) {
		if res.Err != nil {
			return res.Err
		}
		if err := encoding.Encode(res.Upload); err != nil {
			return err
		}
	}
	return nil
}

func export(c *cli.Context, expr activity.Exporter, activityID int64) (*activity.Export, error) {
	log.Info().Int64("activityID", activityID).Msg("export")
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	return expr.Export(ctx, activityID)
}

func qp(c *cli.Context) error {
	expr, err := exporter(c)
	if err != nil {
		return err
	}
	updr, plr, err := uploader(c)
	if err != nil {
		return err
	}
	args := c.Args()
	for i := 0; i < args.Len(); i++ {
		activityID, err := strconv.ParseInt(args.Get(i), 10, 64)
		if err != nil {
			return err
		}
		exp, err := export(c, expr, activityID)
		if err != nil {
			return err
		}
		if err := upload(c, updr, plr, exp); err != nil {
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
