package zwift

import (
	"context"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/commands/activity/internal"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/activity"
	"github.com/bzimmer/gravl/pkg/providers/activity/zwift"
)

const tooSmall = 1024

func NewClient(c *cli.Context) (*zwift.Client, error) {
	client, err := zwift.NewClient(zwift.WithHTTPTracing(c.Bool("http-tracing")))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	username, password := c.String("zwift.username"), c.String("zwift.password")
	token, err := client.Auth.Refresh(ctx, username, password)
	if err != nil {
		return nil, err
	}
	return zwift.NewClient(
		zwift.WithHTTPTracing(c.Bool("http-tracing")),
		zwift.WithToken(token),
	)
}

func athlete(c *cli.Context) error {
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	args := c.Args().Slice()
	if len(args) == 0 {
		args = append(args, zwift.Me)
	}
	for i := range args {
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		profile, err := client.Profile.Profile(ctx, args[i])
		if err != nil {
			return err
		}
		if err = encoding.Encode(profile); err != nil {
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
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	username, password := c.String("zwift.username"), c.String("zwift.password")
	token, err := client.Auth.Refresh(ctx, username, password)
	if err != nil {
		return err
	}
	return encoding.Encode(token)
}

var refreshCommand = &cli.Command{
	Name:   "refresh",
	Usage:  "Acquire a new refresh token",
	Action: refresh,
}

func activities(c *cli.Context) error {
	client, err := NewClient(c)
	if err != nil {
		return err
	}
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
		if err = encoding.Encode(acts[i]); err != nil {
			return err
		}
	}
	return nil
}

var activitiesCommand = &cli.Command{
	Name:    "activities",
	Usage:   "Query activities for an athlete from Strava",
	Aliases: []string{"A"},
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "count",
			Aliases: []string{"N"},
			Value:   0,
			Usage:   "Count",
		},
	},
	Action: activities,
}

func entity(c *cli.Context, f func(context.Context, *zwift.Client, *zwift.Activity) error) error {
	if c.NArg() == 0 {
		log.Warn().Msg("no args specified; exiting")
		return nil
	}
	client, err := NewClient(c)
	if err != nil {
		return err
	}
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
		if err = f(ctx, client, act); err != nil {
			return err
		}
	}
	return nil
}

var activityCommand = &cli.Command{
	Name:    "activity",
	Aliases: []string{"a"},
	Usage:   "Query an activity from Zwift",
	Action: func(c *cli.Context) error {
		return entity(c, func(_ context.Context, _ *zwift.Client, act *zwift.Activity) error {
			return encoding.Encode(act)
		})
	},
}

func export(ctx context.Context, c *cli.Context, client *zwift.Client, act *zwift.Activity) error {
	exp, err := client.Activity.Export(ctx, act)
	if err != nil {
		return err
	}
	if err = internal.Write(c, exp); err != nil {
		return err
	}
	return encoding.Encode(exp)
}

var exportCommand = &cli.Command{
	Name:  "export",
	Usage: "Export a Zwift activity by id",
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
	paths := make([]string, 0)
	if len(args) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			// log but error silently since this is optional behavior
			log.Warn().Err(err).Msg("homedir not found")
			return nil
		}
		// @todo(bzimmer) add windows support when it can be tested
		args = []string{
			filepath.Join(home, "Documents/Zwift/Activities"),
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
			if format != activity.FIT {
				log.Info().Str("file", path).Msg("skipping, not a FIT file")
				return nil
			}
			paths = append(paths, path)
			return nil
		})
		if err != nil {
			return err
		}
	}
	return encoding.Encode(paths)
}

var filesCommand = &cli.Command{
	Name:   "files",
	Usage:  "List all local Zwift files; filters small files (584 bytes) and files named 'inProgressActivity.fit'",
	Action: files,
}

var Command = &cli.Command{
	Name:     "zwift",
	Category: "activity",
	Usage:    "Query Zwift for activities",
	Flags: []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "zwift.username",
			Usage: "Username for the Zwift website",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "zwift.password",
			Usage: "Password for the Zwift website",
		}),
	},
	Subcommands: []*cli.Command{
		activitiesCommand,
		activityCommand,
		athleteCommand,
		exportCommand,
		filesCommand,
		refreshCommand,
	},
}
