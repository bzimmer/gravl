package gravl

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var (
	encoder *json.Encoder
	decoder *json.Decoder
)

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
	// httptracing = c.Bool("http-tracing")
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
			Out:     c.App.ErrWriter,
			NoColor: monochrome,
		},
	)
	return nil
}

func flags() ([]cli.Flag, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
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
			Value: filepath.Join(home, ".gravl.yaml"),
			Usage: "File containing configuration settings",
		},
		&cli.DurationFlag{
			Name:    "timeout",
			Aliases: []string{"t"},
			Value:   time.Millisecond * 10000,
			Usage:   "Timeout duration (eg, 1ms, 2s, 5m, 3h)",
		},
	}, nil
}

func commands() []*cli.Command {
	return []*cli.Command{
		serveCommand,
		wtaCommand,
		gnisCommand,
		visualcrossingCommand,
		tengoCommand,
		rwgpsCommand,
		noaaCommand,
		stravaCommand,
		versionCommand,
	}
}

// Run .
func Run() error {
	fs, err := flags()
	if err != nil {
		return err
	}
	app := &cli.App{
		Name:      "gravl",
		Compiled:  time.Now(),
		HelpName:  "gravl",
		Usage:     "Plan trips",
		UsageText: "gravl - plan trips",
		Flags:     fs,
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
	err = app.RunContext(ctx, os.Args)
	if err != nil {
		return err
	}
	return nil
}
