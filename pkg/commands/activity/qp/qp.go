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

	actcmd "github.com/bzimmer/gravl/pkg/commands/activity"
	"github.com/bzimmer/gravl/pkg/commands/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/pkg/commands/activity/rwgps"
	"github.com/bzimmer/gravl/pkg/commands/activity/strava"
	"github.com/bzimmer/gravl/pkg/commands/activity/zwift"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/activity"
)

type qr struct {
	p activity.Poller
	e activity.Exporter
	u activity.Uploader
	d time.Duration
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
	return &qr{e: e, u: u, p: p, d: d}, nil
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
	var u activity.Uploader
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
		u = client.Uploader()
	case "rwgps", "ridewithgps":
		client, err := rwgps.NewClient(c)
		if err != nil {
			return nil, nil, err
		}
		u = client.Uploader()
	case "strava":
		client, err := strava.NewAPIClient(c)
		if err != nil {
			return nil, nil, err
		}
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
		if err := encoding.Encode(res.Upload); err != nil {
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
	f = append(f, actcmd.RateLimitFlags...)
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
