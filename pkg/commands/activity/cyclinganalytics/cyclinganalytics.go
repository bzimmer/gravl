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

	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/activity/cyclinganalytics"
)

// Maximum number of times to poll status updates on uploads
const polls = 3

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
		files = append(files, &cyclinganalytics.File{Name: base, Reader: file})
		return nil
	})
	return files, err
}

// poll the status possibly following until the operation is completedd
//  https://www.cyclinganalytics.com/developer/api#/user/user_id/upload/upload_id
func poll(ctx context.Context, client *cyclinganalytics.Client, id int64, follow bool) error {
	// status: processing, done, or error
	i, n := 0, polls
	for ; i < n; i++ {
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
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
	if i == n {
		log.Warn().Int("polls", n).Msg("exceeded max polls")
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
				Msg("uploading")
			u, err := client.Rides.Upload(ctx, cyclinganalytics.Me, file)
			if err != nil {
				return err
			}
			if !c.Bool("poll") {
				return encoding.Encode(u)
			}
			return poll(ctx, client, u.UploadID, true)
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
		if err := poll(ctx, client, uploadID, c.Bool("poll")); err != nil {
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
			Name:    "poll",
			Aliases: []string{"p"},
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
		oauthCommand,
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
