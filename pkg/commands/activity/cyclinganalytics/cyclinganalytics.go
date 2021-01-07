package cyclinganalytics

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
)

// Maximum number of times to request status updates on uploads
const maxFollows = 5

func NewClient(c *cli.Context) (*cyclinganalytics.Client, error) {
	return cyclinganalytics.NewClient(
		cyclinganalytics.WithTokenCredentials(
			c.String("cyclinganalytics.access-token"), c.String("cyclinganalytics.refresh-token"), time.Time{}),
		cyclinganalytics.WithAutoRefresh(c.Context),
		cyclinganalytics.WithHTTPTracing(c.Bool("http-tracing")))
}

// collect returns a slice of files for uploading
// Primary use case has been uploading fit files from Zwift so this function
//  filters small files (less then 1K) and files of the name "inProgressActivity.fit"
func collect(name string) ([]*cyclinganalytics.File, error) {
	var files []*cyclinganalytics.File
	err := filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		base := filepath.Base(path)
		if base == "inProgressActivity.fit" {
			log.Warn().
				Str("name", path).
				Msg("skipping, not a completed activity")
			return nil
		}
		if info.Size() < 1024 {
			log.Warn().
				Int64("size", info.Size()).
				Str("name", path).
				Msg("skipping, too small")
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		files = append(files, &cyclinganalytics.File{
			Name:   base,
			Reader: file,
			Size:   info.Size(),
		})
		return nil
	})
	return files, err
}

// follow the status until processing is complete
//  https://www.cyclinganalytics.com/developer/api#/user/user_id/upload/upload_id
func follow(ctx context.Context, client *cyclinganalytics.Client, id int64, follow bool) error {
	// status: processing, done, or error
	i := maxFollows
	for {
		u, err := client.Rides.Status(ctx, cyclinganalytics.Me, id)
		if err != nil {
			return err
		}
		if err = encoding.Encode(u); err != nil {
			return err
		}
		if !(follow && u.Status == "processing") {
			break
		}
		i--
		if i == 0 {
			log.Warn().Int("follows", maxFollows).Msg("exceeded max follows")
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
	return nil
}

func upload(c *cli.Context) error {
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	args := c.Args()
	for i := 0; i < args.Len(); i++ {
		files, err := collect(args.Get(i))
		if err != nil {
			return err
		}
		for _, file := range files {
			defer file.Close()
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			log.Info().
				Str("file", file.Name).
				Int64("size", file.Size).
				Msg("uploading")
			u, err := client.Rides.Upload(ctx, cyclinganalytics.Me, file)
			if err != nil {
				return err
			}
			if c.Bool("follow") {
				return follow(ctx, client, u.UploadID, c.Bool("follow"))
			}
			if err := encoding.Encode(u); err != nil {
				return err
			}
		}
	}
	return nil
}

func status(c *cli.Context) error {
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	args := c.Args()
	for i := 0; i < args.Len(); i++ {
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		uploadID, err := strconv.ParseInt(args.Get(i), 0, 64)
		if err != nil {
			return err
		}
		if err := follow(ctx, client, uploadID, c.Bool("follow")); err != nil {
			return err
		}
	}
	return nil
}

var uploadCommand = &cli.Command{
	Name:    "upload",
	Aliases: []string{"u"},
	Usage:   "Upload an activity file",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "status",
			Aliases: []string{"s"},
			Value:   false,
			Usage:   "Check the status of the upload",
		},
		&cli.BoolFlag{
			Name:    "follow",
			Aliases: []string{"f"},
			Value:   false,
			Usage:   "Continually check the status of the request until it is completed",
		},
	},
	Action: func(c *cli.Context) error {
		if c.Bool("status") {
			return status(c)
		}
		return upload(c)
	},
}

func athlete(c *cli.Context) error {
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	athlete, err := client.User.Me(ctx)
	if err != nil {
		return err
	}
	return encoding.Encode(athlete)
}

var athleteCommand = &cli.Command{
	Name:    "athlete",
	Aliases: []string{"t"},
	Usage:   "Query for the authenticated athlete",
	Action:  athlete,
}

func activities(c *cli.Context) error {
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	rides, err := client.Rides.Rides(ctx, cyclinganalytics.Me)
	if err != nil {
		return err
	}
	for _, ride := range rides {
		err := encoding.Encode(ride)
		if err != nil {
			return err
		}
	}
	return nil
}

var activitiesCommand = &cli.Command{
	Name:    "activities",
	Aliases: []string{"A"},
	Usage:   "Query activities for the authenticated athlete",
	Action:  activities,
}

func activity(c *cli.Context) error {
	client, err := NewClient(c)
	if err != nil {
		return err
	}
	args := c.Args()
	opts := cyclinganalytics.RideOptions{
		Streams: []string{"latitude", "longitude", "elevation"},
	}
	for i := 0; i < args.Len(); i++ {
		ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
		defer cancel()
		rideID, err := strconv.ParseInt(args.Get(i), 0, 64)
		if err != nil {
			return err
		}
		ride, err := client.Rides.Ride(ctx, rideID, opts)
		if err != nil {
			return err
		}
		if err = encoding.Encode(ride); err != nil {
			return err
		}
	}
	return nil
}

var activityCommand = &cli.Command{
	Name:    "activity",
	Aliases: []string{"a"},
	Usage:   "Query an activity for the authenticated athlete",
	Action:  activity,
}

var Command = &cli.Command{
	Name:     "cyclinganalytics",
	Aliases:  []string{"ca"},
	Category: "activity",
	Usage:    "Query CyclingAnalytics",
	Flags:    AuthFlags,
	Subcommands: []*cli.Command{
		activitiesCommand,
		activityCommand,
		athleteCommand,
		uploadCommand,
	},
}

var AuthFlags = []cli.Flag{
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "cyclinganalytics.client-id",
		Usage: "API key for Cycling Analytics API",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "cyclinganalytics.client-secret",
		Usage: "API secret for Cycling Analytics API",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "cyclinganalytics.access-token",
		Usage: "Access token for Cycling Analytics API",
	}),
}
