package hammerhead

import (
	"context"
	"io"
	"os"
	"sync"
	"time"

	api "github.com/bzimmer/activity"
	"github.com/bzimmer/activity/hammerhead"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
	"golang.org/x/time/rate"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/activity"
)

const (
	Provider       = "hammerhead"
	metricActivity = "activity"
)

var (
	before    sync.Once //nolint:gochecknoglobals // once
	errBefore error     //nolint:gochecknoglobals // paired with before
)

func activities(c *cli.Context) error {
	client := gravl.Runtime(c).Hammerhead
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	acts, err := client.Activities.Activities(ctx, api.Pagination{Total: c.Int("count")}, c.String("start-date"))
	if err != nil {
		return err
	}
	enc := gravl.Runtime(c).Encoder
	met := gravl.Runtime(c).Metrics
	met.IncrCounter([]string{Provider, c.Command.Name}, 1)
	for i, act := range acts {
		met.IncrCounter([]string{Provider, metricActivity}, 1)
		log.Info().
			Time("date", act.CreatedAt).
			Str("id", act.ID).
			Str("name", act.Name).
			Msg(c.Command.Name)
		if err = enc.Encode([]any{i, act}); err != nil {
			return err
		}
	}
	return nil
}

func activitiesCommand() *cli.Command {
	return &cli.Command{
		Name:        "activities",
		Aliases:     []string{"A"},
		Usage:       "Query activities for the authenticated athlete",
		Description: "Query the Hammerhead API for a list of activities for the authenticated athlete",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "count",
				Aliases: []string{"N"},
				Value:   0,
				Usage:   "The number of activities to query from Hammerhead",
			},
			&cli.StringFlag{
				Name:  "start-date",
				Usage: "Return activities on or after this date (YYYY-MM-DD)",
			},
		},
		Action: activities,
	}
}

func activityCommand() *cli.Command {
	return &cli.Command{
		Name:        metricActivity,
		Aliases:     []string{"a"},
		Usage:       "Query an activity from Hammerhead",
		Description: "Query the Hammerhead API for a specific activity by its ID",
		Action: func(c *cli.Context) error {
			client := gravl.Runtime(c).Hammerhead
			enc := gravl.Runtime(c).Encoder
			args := c.Args()
			for i := 0; i < args.Len(); i++ {
				err := func() error {
					ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
					defer cancel()
					act, err := client.Activities.Activity(ctx, args.Get(i))
					if err != nil {
						return err
					}
					log.Info().Str("id", act.ID).Str("name", act.Name).Msg(c.Command.Name)
					gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, c.Command.Name}, 1)
					return enc.Encode(act)
				}()
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func writeFile(c *cli.Context, f *api.File) error {
	if f == nil || f.Reader == nil {
		return nil
	}
	if !c.IsSet("overwrite") && !c.IsSet("output") {
		_, err := io.Copy(c.App.Writer, f)
		return err
	}
	filename := f.Filename
	if c.IsSet("output") {
		filename = c.String("output")
	}
	fs := gravl.Runtime(c).Fs
	if _, err := fs.Stat(filename); err == nil && !c.Bool("overwrite") {
		log.Error().Str("filename", filename).Msg("file exists and -o flag not specified")
		return os.ErrExist
	}
	fp, err := fs.Create(filename)
	if err != nil {
		return err
	}
	defer func() { _ = fp.(afero.File).Close() }()
	_, err = io.Copy(fp, f)
	if err != nil {
		return err
	}
	return gravl.Runtime(c).Encoder.Encode(map[string]string{
		"filename": filename,
		"format":   f.Format.String(),
	})
}

func fileCommand() *cli.Command {
	return &cli.Command{
		Name:        "file",
		Aliases:     []string{"f"},
		Usage:       "Download a FIT file for an activity from Hammerhead",
		Description: "Download the original FIT file for a specific Hammerhead activity by its ID; streams to stdout if --output is not set",
		ArgsUsage:   "ACTIVITY_ID (...)",
		Flags: []cli.Flag{
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
				Usage:   "The filename to write the FIT file to; if not specified the contents are streamed to stdout",
			},
		},
		Action: func(c *cli.Context) error {
			client := gravl.Runtime(c).Hammerhead
			args := c.Args()
			for i := 0; i < args.Len(); i++ {
				err := func() error {
					ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
					defer cancel()
					id := args.Get(i)
					f, err := client.Activities.File(ctx, id)
					if err != nil {
						return err
					}
					defer f.Close()
					log.Info().Str("id", id).Str("filename", f.Filename).Msg(c.Command.Name)
					gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, c.Command.Name}, 1)
					return writeFile(c, f)
				}()
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func oauthCommand() *cli.Command {
	return activity.OAuthCommand(&activity.OAuthConfig{
		Port:     9004,
		Provider: Provider,
		Scopes:   []string{"activity:read"},
	})
}

func Before(c *cli.Context) error {
	before.Do(func() {
		var client *hammerhead.Client
		client, errBefore = hammerhead.NewClient(
			hammerhead.WithClientCredentials(c.String("hammerhead-client-id"), c.String("hammerhead-client-secret")),
			hammerhead.WithTokenCredentials(
				c.String("hammerhead-access-token"), c.String("hammerhead-refresh-token"), time.Now().Add(-1*time.Minute)),
			hammerhead.WithAutoRefresh(c.Context),
			hammerhead.WithHTTPTracing(c.Bool("http-tracing")),
			hammerhead.WithRateLimiter(rate.NewLimiter(
				rate.Every(c.Duration("rate-limit")), c.Int("rate-burst"))))
		if errBefore != nil {
			return
		}
		gravl.Runtime(c).Endpoints[Provider] = hammerhead.Endpoint()
		gravl.Runtime(c).Hammerhead = client
		gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, "client", "created"}, 1)
		log.Info().Msg("created hammerhead client")
	})
	return errBefore
}

func Command() *cli.Command {
	return &cli.Command{
		Name:        Provider,
		Category:    metricActivity,
		Usage:       "Query Hammerhead for rides and activities",
		Description: "Operations supported by the Hammerhead Karoo API",
		Flags:       append(AuthFlags(), activity.RateLimitFlags()...),
		Before:      Before,
		Subcommands: []*cli.Command{
			activitiesCommand(),
			activityCommand(),
			fileCommand(),
			oauthCommand(),
		},
	}
}

func AuthFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "hammerhead-client-id",
			Usage:   "Hammerhead client id",
			EnvVars: []string{"HAMMERHEAD_CLIENT_ID"},
		},
		&cli.StringFlag{
			Name:    "hammerhead-client-secret",
			Usage:   "Hammerhead client secret",
			EnvVars: []string{"HAMMERHEAD_CLIENT_SECRET"},
		},
		&cli.StringFlag{
			Name:    "hammerhead-access-token",
			Usage:   "Hammerhead access token",
			EnvVars: []string{"HAMMERHEAD_ACCESS_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "hammerhead-refresh-token",
			Usage:   "Hammerhead refresh token",
			EnvVars: []string{"HAMMERHEAD_REFRESH_TOKEN"},
		},
	}
}
