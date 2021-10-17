package qp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"github.com/bzimmer/activity"
	"github.com/bzimmer/gravl/pkg"
	actcmd "github.com/bzimmer/gravl/pkg/activity"
	"github.com/bzimmer/gravl/pkg/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/pkg/activity/rwgps"
	"github.com/bzimmer/gravl/pkg/activity/strava"
	"github.com/bzimmer/gravl/pkg/activity/zwift"
)

type qr struct {
	p   activity.Poller
	e   activity.Exporter
	u   activity.Uploader
	d   time.Duration
	enc pkg.Encoder
}

func newqr(c *cli.Context) (*qr, error) {
	e, err := exporter(c)
	if err != nil {
		return nil, err
	}
	u, p, err := uploader(c)
	if err != nil {
		return nil, err
	}
	d := c.Duration("timeout")
	enc := pkg.Runtime(c).Encoder
	return &qr{e: e, u: u, p: p, d: d, enc: enc}, nil
}

func exporter(c *cli.Context) (activity.Exporter, error) {
	exp := c.String("exporter")
	log.Info().Str("provider", exp).Msg("exporter")
	switch strings.ToLower(exp) {
	case "":
		return nil, errors.New("exporter is a required flag")
	case "zwift":
		client := pkg.Runtime(c).Zwift
		return client.Exporter(), nil
	case "strava": // nolint
		return nil, errors.New("strava exporter needs to be rewritten to use gpx file")
	}
	return nil, fmt.Errorf("unknown exporter {%s}", exp)
}

func uploader(c *cli.Context) (activity.Uploader, activity.Poller, error) {
	var u activity.Uploader
	upd := c.String("uploader")
	log.Info().Str("provider", upd).Msg("uploader")
	switch strings.ToLower(upd) {
	case "":
		return nil, nil, errors.New("uploader is a required flag")
	case "ca", "cyclinganalytics":
		client := pkg.Runtime(c).CyclingAnalytics
		u = client.Uploader()
	case "rwgps", "ridewithgps":
		client := pkg.Runtime(c).RideWithGPS
		u = client.Uploader()
	case "strava":
		client := pkg.Runtime(c).Strava
		u = client.Uploader()
	default:
		return nil, nil, fmt.Errorf("unknown uploader {%s}", upd)
	}
	p := activity.NewPoller(u,
		activity.WithInterval(c.Duration("interval")),
		activity.WithIterations(c.Int("iterations")))

	return u, p, nil
}

func (q *qr) upload(ctx context.Context, export *activity.Export) error {
	log.Info().Int64("activityID", export.ID).Msg("upload")
	out := new(bytes.Buffer)
	_, err := io.Copy(out, export)
	if err != nil {
		return err
	}
	file := &activity.File{Reader: out, Format: export.Format, Name: export.Name}
	defer file.Close()
	ctx, cancel := context.WithTimeout(ctx, q.d)
	defer cancel()
	u, err := q.u.Upload(ctx, file)
	if err != nil {
		return err
	}
	for res := range q.p.Poll(ctx, u.Identifier()) {
		if res.Err != nil {
			return res.Err
		}
		if err := q.enc.Encode(res.Upload); err != nil {
			return err
		}
	}
	return ctx.Err()
}

func (q *qr) export(ctx context.Context, activityID int64) (*activity.Export, error) {
	log.Info().Int64("activityID", activityID).Msg("export")
	ctx, cancel := context.WithTimeout(ctx, q.d)
	defer cancel()
	return q.e.Export(ctx, activityID)
}

func qp(c *cli.Context) error {
	q, err := newqr(c)
	if err != nil {
		return err
	}
	args := c.Args()
	grp, ctx := errgroup.WithContext(c.Context)
	for i := 0; i < args.Len(); i++ {
		activityID, err := strconv.ParseInt(args.Get(i), 10, 64)
		if err != nil {
			return err
		}
		grp.Go(func() error {
			exp, err := q.export(ctx, activityID)
			if err != nil {
				return err
			}
			if err := q.upload(ctx, exp); err != nil {
				return err
			}
			return nil
		})
	}
	return grp.Wait()
}

var flags = func() []cli.Flag {
	f := []cli.Flag{
		&cli.StringFlag{
			Name:    "exporter",
			Aliases: []string{"e"},
			Usage:   "Export data provider"},
		&cli.StringFlag{
			Name:    "uploader",
			Aliases: []string{"u"},
			Usage:   "Upload data provider"},
	}
	for _, x := range [][]cli.Flag{
		actcmd.RateLimitFlags, cyclinganalytics.AuthFlags, rwgps.AuthFlags, strava.AuthFlags, zwift.AuthFlags,
	} {
		f = append(f, x...)
	}
	return f
}()

var Command = &cli.Command{
	Name:      "qp",
	Category:  "activity",
	Usage:     "Copy an activity from an exporter to an uploader",
	ArgsUsage: "ACTIVITY_ID (...)",
	Flags:     flags,
	Action:    qp,
	Before:    pkg.Befores(strava.Before, cyclinganalytics.Before, zwift.Before, rwgps.Before),
}
