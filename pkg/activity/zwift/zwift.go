package zwift

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
	"golang.org/x/time/rate"

	api "github.com/bzimmer/activity"
	"github.com/bzimmer/activity/zwift"
	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/activity"
)

const (
	tooSmall = 1024
	Provider = "zwift"
)

var before sync.Once

func athlete(c *cli.Context) error {
	client := pkg.Runtime(c).Zwift
	args := c.Args().Slice()
	if len(args) == 0 {
		args = append(args, zwift.Me)
	}
	enc := pkg.Runtime(c).Encoder
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	for i := range args {
		profile, err := client.Profile.Profile(ctx, args[i])
		if err != nil {
			return err
		}
		log.Info().Int64("id", profile.ID).Str("username", profile.PublicID).Msg(c.Command.Name)
		pkg.Runtime(c).Metrics.IncrCounter([]string{Provider, c.Command.Name}, 1)
		if err := enc.Encode(profile); err != nil {
			return err
		}
	}
	return nil
}

func athleteCommand() *cli.Command {
	return &cli.Command{
		Name:    "athlete",
		Usage:   "Query the athlete profile from Zwift",
		Aliases: []string{"t"},
		Action:  athlete,
	}
}

func refresh(c *cli.Context) error {
	client := pkg.Runtime(c).Zwift
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	username, password := c.String("zwift.username"), c.String("zwift.password")
	token, err := client.Auth.Refresh(ctx, username, password)
	if err != nil {
		return err
	}
	return pkg.Runtime(c).Encoder.Encode(token)
}

func refreshCommand() *cli.Command {
	return &cli.Command{
		Name:   "refresh",
		Usage:  "Acquire a new refresh token",
		Action: refresh,
	}
}

func activities(c *cli.Context) error {
	client := pkg.Runtime(c).Zwift
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	profile, err := client.Profile.Profile(ctx, zwift.Me)
	if err != nil {
		return err
	}
	spec := api.Pagination{Total: c.Int("count")}
	acts, err := client.Activity.Activities(ctx, profile.ID, spec)
	if err != nil {
		return err
	}
	pkg.Runtime(c).Metrics.IncrCounter([]string{Provider, c.Command.Name}, 1)
	for _, act := range acts {
		pkg.Runtime(c).Metrics.IncrCounter([]string{Provider, "activity"}, 1)
		log.Info().
			Time("date", act.StartDate.Time).
			Int64("id", act.ID).
			Str("name", act.Name).
			Msg(c.Command.Name)
		if err := pkg.Runtime(c).Encoder.Encode(act); err != nil {
			return err
		}
	}
	return nil
}

func activitiesCommand() *cli.Command {
	return &cli.Command{
		Name:    "activities",
		Usage:   "Query activities for an athlete from Zwift",
		Aliases: []string{"A"},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "count",
				Aliases: []string{"N"},
				Value:   0,
				Usage:   "The number of activities to query from Zwift (the number returned will be <= N)",
			},
		},
		Action: activities,
	}
}

func entity(c *cli.Context, f func(context.Context, *zwift.Activity) error) error {
	if c.NArg() == 0 {
		log.Warn().Msg("no args specified; exiting")
		return nil
	}
	client := pkg.Runtime(c).Zwift
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	profile, err := client.Profile.Profile(ctx, zwift.Me)
	if err != nil {
		return err
	}
	for _, arg := range c.Args().Slice() {
		x, err := strconv.ParseInt(arg, 0, 64)
		if err != nil {
			return err
		}
		log.Info().Int64("id", x).Str("entity", c.Command.Name).Msg("querying")
		pkg.Runtime(c).Metrics.IncrCounter([]string{Provider, c.Command.Name}, 1)
		act, err := client.Activity.Activity(ctx, profile.ID, x)
		if err != nil {
			return err
		}
		if err := f(ctx, act); err != nil {
			return err
		}
	}
	return nil
}

func activityCommand() *cli.Command {
	return &cli.Command{
		Name:      "activity",
		Aliases:   []string{"a"},
		Usage:     "Query an activity from Zwift",
		ArgsUsage: "ACTIVITY_ID (...)",
		Action: func(c *cli.Context) error {
			return entity(c, func(_ context.Context, act *zwift.Activity) error {
				return pkg.Runtime(c).Encoder.Encode(act)
			})
		},
	}
}

// Primary use case has been uploading fit files from a local Zwift directory
// Filters small files (584 bytes) and files named "inProgressActivity.fit"
// If no arguments are specified will try to default to the Zwift Activities directory
func files(c *cli.Context) error {
	args := c.Args().Slice()
	if len(args) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			// log but error silently since this is optional behavior
			log.Warn().Err(err).Msg("homedir not found")
			return nil
		}
		args = append(args, filepath.Join(home, "Documents", "Zwift", "Activities"))
	}
	fs := pkg.Runtime(c).Fs
	enc := pkg.Runtime(c).Encoder
	met := pkg.Runtime(c).Metrics
	log.Info().Str("fs", fs.Name()).Msg("walk")
	for _, arg := range args {
		err := afero.Walk(fs, arg, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if os.IsNotExist(err) {
					met.IncrCounter([]string{Provider, c.Command.Name, "skipping", "does-not-exist"}, 1)
					log.Warn().Str("file", path).Msg("path does not exist")
					return nil
				}
				return err
			}
			met.IncrCounter([]string{Provider, c.Command.Name, "found"}, 1)
			if info.IsDir() {
				met.IncrCounter([]string{Provider, c.Command.Name, "directory"}, 1)
				return nil
			}
			base := filepath.Base(path)
			if base == "inProgressActivity.fit" {
				met.IncrCounter([]string{Provider, c.Command.Name, "skipping", "in-progress"}, 1)
				log.Warn().Str("file", path).Msg("skipping, activity in progress")
				return nil
			}
			if info.Size() <= tooSmall {
				met.IncrCounter([]string{Provider, c.Command.Name, "skipping", "too-small"}, 1)
				log.Warn().Int64("size", info.Size()).Str("file", path).Msg("skipping, too small")
				return nil
			}
			format := api.ToFormat(filepath.Ext(path))
			if format != api.FormatFIT {
				met.IncrCounter([]string{Provider, c.Command.Name, "skipping", format.String()}, 1)
				log.Info().Str("file", path).Msg("skipping, not a FIT file")
				return nil
			}
			met.IncrCounter([]string{Provider, c.Command.Name, "success"}, 1)
			return enc.Encode(path)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func filesCommand() *cli.Command {
	return &cli.Command{
		Name:   "files",
		Usage:  "List all local Zwift files",
		Action: files,
	}
}

func AuthFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "zwift-username",
			Usage:   "Zwift username",
			EnvVars: []string{"ZWIFT_USERNAME"},
		},
		&cli.StringFlag{
			Name:    "zwift-password",
			Usage:   "Zwift password",
			EnvVars: []string{"ZWIFT_PASSWORD"},
		},
	}
}

// Before configures the zwift client
func Before(c *cli.Context) error {
	var err error
	before.Do(func() {
		var client *zwift.Client
		client, err = zwift.NewClient(
			zwift.WithTokenRefresh(c.String("zwift-username"), c.String("zwift-password")),
			zwift.WithHTTPTracing(c.Bool("http-tracing")),
			zwift.WithRateLimiter(rate.NewLimiter(
				rate.Every(c.Duration("rate-limit")), c.Int("rate-burst"))))
		if err != nil {
			return
		}
		pkg.Runtime(c).Endpoints[Provider] = zwift.Endpoint()
		pkg.Runtime(c).Zwift = client
		log.Info().Msg("created zwift client")
	})
	return err
}

func Command() *cli.Command {
	return &cli.Command{
		Name:        "zwift",
		Category:    "activity",
		Usage:       "Query Zwift for activities",
		Description: "Operations supported by the Zwift API",
		Flags:       append(AuthFlags(), activity.RateLimitFlags()...),
		Before:      Before,
		Subcommands: []*cli.Command{
			activitiesCommand(),
			activityCommand(),
			athleteCommand(),
			filesCommand(),
			refreshCommand(),
		},
	}
}
