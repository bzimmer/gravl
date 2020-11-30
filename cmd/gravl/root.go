package gravl

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	stdlog "log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/adrg/xdg"
	"github.com/bzimmer/gravl/pkg"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var (
	encoder *json.Encoder
	decoder *json.Decoder
)

type logger struct{}

func (w logger) Write(p []byte) (n int, err error) {
	s := strings.TrimSuffix(string(p), "\n")
	log.Debug().Msg(s)
	return len(p), nil
}

func mustRandomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

func merge(flags ...[]cli.Flag) []cli.Flag {
	f := make([]cli.Flag, 0)
	for _, x := range flags {
		f = append(f, x...)
	}
	return f
}

func before(befores ...cli.BeforeFunc) cli.BeforeFunc {
	return func(c *cli.Context) error {
		for _, b := range befores {
			if b == nil {
				continue
			}
			err := b(c)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func initConfig(c *cli.Context) error {
	config := func() (altsrc.InputSourceContext, error) {
		return altsrc.NewYamlSourceFromFile(c.String("config"))
	}
	for _, cmd := range c.App.Commands {
		cmd.Before = before(
			altsrc.InitInputSource(cmd.Flags, config),
			cmd.Before)
	}
	return nil
}

func initFlags(c *cli.Context) error {
	return nil
}

func initEncoding(c *cli.Context) error {
	decoder = newDecoder(c.App.Reader)
	encoder = newEncoder(c.App.Writer, c.Bool("compact"))
	return nil
}

func initLogging(c *cli.Context) error {
	monochrome := c.Bool("monochrome")
	level, err := zerolog.ParseLevel(c.String("verbosity"))
	if err != nil {
		return err
	}
	color.NoColor = monochrome
	zerolog.SetGlobalLevel(level)
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.DurationFieldInteger = true
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:        c.App.ErrWriter,
			NoColor:    monochrome,
			TimeFormat: time.RFC3339,
		},
	)
	stdlog.SetOutput(logger{})

	return nil
}

func flags() []cli.Flag {
	config := path.Join(xdg.ConfigHome, pkg.PackageName, "gravl.yaml")
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "verbosity",
			Aliases: []string{"v"},
			Value:   "info",
			Usage:   "Log level (trace, debug, info, warn, error, fatal, panic)",
		},
		&cli.BoolFlag{
			Name:    "monochrome",
			Aliases: []string{"m"},
			Value:   false,
			Usage:   "Use monochrome logging, color enabled by default",
		},
		&cli.BoolFlag{
			Name:    "compact",
			Aliases: []string{"c"},
			Value:   false,
			Usage:   "Use compact JSON output",
		},
		&cli.BoolFlag{
			Name:  "http-tracing",
			Value: false,
			Usage: "Log all http calls (warning: this will log ids, keys, and other sensitive information)",
		},
		&cli.StringFlag{
			Name:  "config",
			Value: config,
			Usage: "File containing configuration settings",
		},
		&cli.DurationFlag{
			Name:    "timeout",
			Aliases: []string{"t"},
			Value:   time.Millisecond * 10000,
			Usage:   "Timeout duration (eg, 1ms, 2s, 5m, 3h)",
		},
	}
}

func commands() []*cli.Command {
	return []*cli.Command{
		cyclinganalyticsCommand,
		gnisCommand,
		noaaCommand,
		openweatherCommand,
		rwgpsCommand,
		serveCommand,
		srtmCommand,
		stravaCommand,
		versionCommand,
		visualcrossingCommand,
		wtaCommand,
	}
}

// Run .
func Run() error {
	app := &cli.App{
		Name:      "gravl",
		HelpName:  "gravl",
		Usage:     "Plan trips",
		UsageText: "gravl - plan trips",
		Flags:     flags(),
		Commands:  commands(),
		Before: func(c *cli.Context) error {
			fns := []cli.BeforeFunc{initFlags, initLogging, initEncoding, initConfig}
			for _, fn := range fns {
				if e := fn(c); e != nil {
					return e
				}
			}
			return nil
		},
		ExitErrHandler: func(c *cli.Context, err error) {
			log.Error().Err(err).Msg("gravl")
		},
	}
	ctx := context.Background()
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		return err
	}
	return nil
}
