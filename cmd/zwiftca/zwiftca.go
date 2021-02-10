package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/commands/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/pkg/commands/activity/strava"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/commands/gravl"
	"github.com/bzimmer/gravl/pkg/providers/activity"
	caapi "github.com/bzimmer/gravl/pkg/providers/activity/cyclinganalytics"
	stravaapi "github.com/bzimmer/gravl/pkg/providers/activity/strava"
	stravaweb "github.com/bzimmer/gravl/pkg/providers/activity/strava/web"
)

const polls = 3

func list(c *cli.Context) error {
	client, err := strava.NewAPIClient(c)
	if err != nil {
		return err
	}
	var ok bool
	var act *stravaapi.Activity
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	acts, errs := client.Activity.Activities(ctx, activity.Pagination{Total: c.Int("count")})
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err, ok = <-errs:
			// if the channel was not closed an error occurred so return it
			// if the channel is closed do nothing to ensure the activity channel can run to
			//  completion and return the full slice of activities
			if ok {
				return err
			}
		case act, ok = <-acts:
			if !ok {
				// the channel is closed, done
				return nil
			}
			if act.Type == "VirtualRide" {
				m := []interface{}{act.ID, act.Name, act.StartDate}
				if err = encoding.Encode(m); err != nil {
					return err
				}
			}
		}
	}
}

var listCommand = &cli.Command{
	Name:  "list",
	Usage: "List VirtualRide activities",
	Flags: []cli.Flag{
		&cli.IntFlag{Name: "count", Aliases: []string{"N"}, Value: 100, Usage: "Count"},
	},
	Action: list,
}

func export(ctx context.Context, c *cli.Context, activityID int64) (*stravaweb.Export, error) {
	client, err := strava.NewWebClient(c)
	if err != nil {
		return nil, err
	}
	log.Info().Int64("activityID", activityID).Msg("exporting")
	return client.Export.Export(ctx, activityID, stravaweb.Original)
}

func upload(ctx context.Context, c *cli.Context, export *stravaweb.Export) (*caapi.Upload, error) {
	client, err := cyclinganalytics.NewClient(c)
	if err != nil {
		return nil, err
	}
	out := new(bytes.Buffer)
	_, err = io.Copy(out, export)
	if err != nil {
		return nil, err
	}
	file := &caapi.File{Reader: out, Name: export.Name}
	defer file.Close()
	log.Info().Str("name", export.Name).Msg("uploading")
	return client.Rides.Upload(ctx, caapi.Me, file)
}

func sync(c *cli.Context) error {
	client, err := cyclinganalytics.NewClient(c)
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
		exp, err := export(ctx, c, actID)
		if err != nil {
			return err
		}
		sts, err := upload(ctx, c, exp)
		if err != nil {
			return err
		}
		if err = encoding.Encode(sts); err != nil {
			return err
		}
		i, n := 0, polls
		for ; i < n && sts.Status == "processing"; i++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(2 * time.Second):
				sts, err = client.Rides.Status(ctx, sts.UserID, sts.UploadID)
				if err != nil {
					return err
				}
				if err = encoding.Encode(sts); err != nil {
					return err
				}
			}
		}
		if i == n {
			log.Warn().Int("polls", n).Msg("exceeded max polls")
		}
	}
	return nil
}

var syncCommand = &cli.Command{
	Name:   "sync",
	Usage:  "Sync the VirtualRide activity in Strava to CyclingAnalytics",
	Action: sync,
}

func main() {
	flags := append(gravl.Flags, gravl.ConfigFlag("gravl.yaml"))
	flags = append(flags, strava.AuthFlags...)
	flags = append(flags, cyclinganalytics.AuthFlags...)
	app := &cli.App{
		Name:     "zwiftca",
		HelpName: "zwiftca",
		Flags:    flags,
		Before:   gravl.Befores(gravl.InitLogging(), gravl.InitEncoding(), gravl.InitConfig()),
		Commands: []*cli.Command{listCommand, syncCommand},
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
