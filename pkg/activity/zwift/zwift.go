package zwift

import (
	"context"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/time/rate"

	"github.com/bzimmer/activity"
	"github.com/bzimmer/activity/zwift"
	"github.com/bzimmer/gravl/pkg"
	actcmd "github.com/bzimmer/gravl/pkg/activity"
)

const (
	tooSmall = 1024
	provider = "zwift"
)

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
		pkg.Runtime(c).Metrics.IncrCounter([]string{provider, c.Command.Name}, 1)
		if err := enc.Encode(profile); err != nil {
			return err
		}
	}
	return nil
}

var athleteCommand = &cli.Command{
	Name:    "athlete",
	Usage:   "Query the athlete profile from Zwift",
	Aliases: []string{"t"},
	Action:  athlete,
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

var refreshCommand = &cli.Command{
	Name:   "refresh",
	Usage:  "Acquire a new refresh token",
	Action: refresh,
}

func activities(c *cli.Context) error {
	client := pkg.Runtime(c).Zwift
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	profile, err := client.Profile.Profile(ctx, zwift.Me)
	if err != nil {
		return err
	}
	spec := activity.Pagination{Total: c.Int("count")}
	acts, err := client.Activity.Activities(ctx, profile.ID, spec)
	if err != nil {
		return err
	}
	for i := range acts {
		if err := pkg.Runtime(c).Encoder.Encode(acts[i]); err != nil {
			return err
		}
	}
	return nil
}

var activitiesCommand = &cli.Command{
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

func entity(c *cli.Context, f func(context.Context, *zwift.Client, *zwift.Activity) error) error {
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
		log.Info().Int64("id", x).Msg("querying activity")
		act, err := client.Activity.Activity(ctx, profile.ID, x)
		if err != nil {
			return err
		}
		if err := f(ctx, client, act); err != nil {
			return err
		}
	}
	return nil
}

var activityCommand = &cli.Command{
	Name:      "activity",
	Aliases:   []string{"a"},
	Usage:     "Query an activity from Zwift",
	ArgsUsage: "ACTIVITY_ID (...)",
	Action: func(c *cli.Context) error {
		return entity(c, func(_ context.Context, _ *zwift.Client, act *zwift.Activity) error {
			return pkg.Runtime(c).Encoder.Encode(act)
		})
	},
}

func export(ctx context.Context, c *cli.Context, client *zwift.Client, act *zwift.Activity) error {
	exp, err := client.Activity.ExportActivity(ctx, act)
	if err != nil {
		return err
	}
	return actcmd.Write(c, exp)
}

var exportCommand = &cli.Command{
	Name:      "export",
	Usage:     "Export a Zwift activity by id",
	ArgsUsage: "ACTIVITY_ID (...)",
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
			Usage:   "The filename to use for writing the contents of the export, if not specified the contents are streamed to stdout",
		},
	},
	Action: func(c *cli.Context) error {
		return entity(c, func(ctx context.Context, client *zwift.Client, act *zwift.Activity) error {
			return export(ctx, c, client, act)
		})
	},
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
		// @todo(bzimmer) add windows support when it can be tested
		args = []string{
			filepath.Join(home, "Documents", "Zwift", "Activities"),
		}
	}
	for _, arg := range args {
		err := filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				// log and continue
				log.Warn().Err(err).Str("file", path).Msg("path does not exist")
				return nil
			}
			if info.IsDir() {
				return nil
			}
			base := filepath.Base(path)
			if base == "inProgressActivity.fit" {
				log.Warn().Str("file", path).Msg("skipping, activity in progress")
				return nil
			}
			if info.Size() <= tooSmall {
				log.Warn().Int64("size", info.Size()).Str("file", path).Msg("skipping, too small")
				return nil
			}
			format := activity.ToFormat(filepath.Ext(path))
			if format != activity.FormatFIT {
				log.Info().Str("file", path).Msg("skipping, not a FIT file")
				return nil
			}
			return pkg.Runtime(c).Encoder.Encode(path)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

var filesCommand = &cli.Command{
	Name:   "files",
	Usage:  "List all local Zwift files",
	Action: files,
}

var AuthFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "zwift-username",
		Usage:   "zwift username",
		EnvVars: []string{"ZWIFT_USERNAME"},
	},
	&cli.StringFlag{
		Name:    "zwift-password",
		Usage:   "zwift password",
		EnvVars: []string{"ZWIFT_PASSWORD"},
	},
}

func Before(c *cli.Context) error {
	client, err := zwift.NewClient(zwift.WithHTTPTracing(c.Bool("http-tracing")))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	username, password := c.String("zwift-username"), c.String("zwift-password")
	token, err := client.Auth.Refresh(ctx, username, password)
	if err != nil {
		return err
	}
	client, err = zwift.NewClient(
		zwift.WithHTTPTracing(c.Bool("http-tracing")),
		zwift.WithToken(token),
		zwift.WithRateLimiter(rate.NewLimiter(
			rate.Every(c.Duration("rate-limit")), c.Int("rate-burst"))))
	if err != nil {
		return err
	}
	pkg.Runtime(c).Zwift = client
	return nil
}

func Command() *cli.Command {
	return &cli.Command{
		Name:        "zwift",
		Category:    "activity",
		Usage:       "Query Zwift for activities",
		Description: "Operations supported by the Zwift API",
		Flags:       append(AuthFlags, actcmd.RateLimitFlags...),
		Before:      Before,
		Subcommands: []*cli.Command{
			activitiesCommand,
			activityCommand,
			athleteCommand,
			exportCommand,
			filesCommand,
			refreshCommand,
		},
	}
}
