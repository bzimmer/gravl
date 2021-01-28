package strava

import (
	"bytes"
	"context"
	"html/template"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/time/rate"

	"github.com/bzimmer/gravl/pkg/commands/encoding"
	stravaweb "github.com/bzimmer/gravl/pkg/providers/activity/strava/web"
)

func NewWebClient(c *cli.Context) (*stravaweb.Client, error) {
	client, err := stravaweb.NewClient(
		stravaweb.WithHTTPTracing(c.Bool("http-tracing")),
		stravaweb.WithCookieJar(),
		stravaweb.WithRateLimiter(rate.NewLimiter(rate.Every(2*time.Second), 5)))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	username, password := c.String("strava.username"), c.String("strava.password")
	if err = client.Auth.Login(ctx, username, password); err != nil {
		return nil, err
	}
	return client, nil
}

func export(c *cli.Context) error {
	client, err := NewWebClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	format := stravaweb.ToFormat(c.String("format"))
	args := c.Args().Slice()
	for i := 0; i < len(args); i++ {
		x, err := strconv.ParseInt(args[i], 0, 64)
		if err != nil {
			return err
		}
		var file *stravaweb.ExportFile
		file, err = client.Export.Export(ctx, x, format)
		if err != nil {
			return err
		}
		fn := file.Name
		if c.IsSet("template") {
			var t *template.Template
			t, err = template.New("export").Parse(c.String("template"))
			if err != nil {
				return err
			}
			var out bytes.Buffer
			err = t.Execute(&out, file)
			if err != nil {
				return err
			}
			fn = out.String()
		}
		if _, err = os.Stat(fn); err == nil && !c.Bool("overwrite") {
			log.Error().Str("filename", fn).Msg("file exists and -o flag not specified")
			return os.ErrExist
		}
		out, err := os.Create(fn)
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, file.Reader)
		if err != nil {
			return err
		}
		file.Name = fn
		if err = encoding.Encode(file); err != nil {
			return err
		}
	}
	return nil
}

var exportCommand = &cli.Command{
	Name:  "export",
	Usage: "export a Strava activity by id, optionally specifying the format and filename template",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"F"},
			Value:   stravaweb.Original.String(),
			Usage:   "Export data file in the specified format",
		},
		&cli.StringFlag{
			Name:    "template",
			Aliases: []string{"T"},
			Usage:   "Export data filename template; fields: ID, Name, Format, Extension",
		},
		&cli.BoolFlag{
			Name:    "overwrite",
			Aliases: []string{"o"},
			Value:   false,
			Usage:   "Overwrite the file if it exists; fail otherwise",
		},
	},
	Action: export,
}

func fitness(c *cli.Context) error {
	webclient, err := NewWebClient(c)
	if err != nil {
		return err
	}
	apiclient, err := NewAPIClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	athlete, err := apiclient.Athlete.Athlete(ctx)
	if err != nil {
		return err
	}
	load, err := webclient.Fitness.TrainingLoad(ctx, athlete.ID)
	if err != nil {
		return err
	}
	return encoding.Encode(load)
}

var fitnessCommand = &cli.Command{
	Name:   "fitness",
	Usage:  "Query Strava for training load data",
	Action: fitness,
}
