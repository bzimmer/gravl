package gravl

import (
	"context"
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/xdg"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/commands"
	"github.com/bzimmer/gravl/pkg/commands/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/pkg/commands/activity/rwgps"
	"github.com/bzimmer/gravl/pkg/commands/activity/strava"
	"github.com/bzimmer/gravl/pkg/commands/activity/wta"
	"github.com/bzimmer/gravl/pkg/commands/analysis"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/commands/geo/gnis"
	"github.com/bzimmer/gravl/pkg/commands/geo/gpx"
	"github.com/bzimmer/gravl/pkg/commands/geo/srtm"
	"github.com/bzimmer/gravl/pkg/commands/store"
	"github.com/bzimmer/gravl/pkg/commands/version"
	"github.com/bzimmer/gravl/pkg/commands/wx/noaa"
	"github.com/bzimmer/gravl/pkg/commands/wx/openweather"
	"github.com/bzimmer/gravl/pkg/commands/wx/visualcrossing"
)

type logger struct{}

func (w logger) Write(p []byte) (n int, err error) {
	s := strings.TrimSuffix(string(p), "\n")
	log.Debug().Msg(s)
	return len(p), nil
}

func flatten(cmds []*cli.Command) []*cli.Command {
	var res []*cli.Command
	for i := range cmds {
		res = append(res, cmds[i])
		res = append(res, flatten(cmds[i].Subcommands)...)
	}
	return res
}

func initConfig(c *cli.Context) error {
	cfg := c.String("config")
	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		log.Error().
			Str("path", cfg).
			Msg("unable to find config file")
		return errors.New("invalid config file")
	}
	config := func() (altsrc.InputSourceContext, error) {
		return altsrc.NewYamlSourceFromFile(cfg)
	}
	for _, cmd := range flatten(c.App.Commands) {
		cmd.Before = commands.Before(altsrc.InitInputSource(cmd.Flags, config), cmd.Before)
	}
	return nil
}

func initEncoding(c *cli.Context) error {
	encoder, err := encoding.NewEncoder(c.App.Writer, c.String("encoding"), c.Bool("compact"))
	if err != nil {
		return err
	}
	encoding.Encode = encoder.Encode
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
	zerolog.DurationFieldInteger = false
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

var commandsCommand = &cli.Command{
	Name:  "commands",
	Usage: "Return all possible commands",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "relative",
			Aliases: []string{"r"},
			Usage:   "Specify the command relative to the current working directory",
		},
	},
	Action: func(c *cli.Context) error {
		var commands []string
		var commander func(string, []*cli.Command)
		commander = func(prefix string, cmds []*cli.Command) {
			for i := range cmds {
				cmd := fmt.Sprintf("%s %s", prefix, cmds[i].Name)
				if !cmds[i].Hidden && cmds[i].Action != nil {
					commands = append(commands, cmd)
				}
				commander(cmd, cmds[i].Subcommands)
			}
		}
		cmd := c.App.Name
		if c.Bool("relative") {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			cmd, err = os.Executable()
			if err != nil {
				return err
			}
			cmd, err = filepath.Rel(cwd, cmd)
			if err != nil {
				return err
			}
		}
		commander(cmd, c.App.Commands)
		return encoding.Encode(commands)
	},
}

// ConfigFlag for the default gravl configuration file
var ConfigFlag = func() cli.Flag {
	config := path.Join(xdg.ConfigHome, pkg.PackageName, "gravl.yaml")
	return &cli.PathFlag{
		Name:  "config",
		Value: config,
		Usage: "File containing configuration settings",
	}
}()

var flags = func() []cli.Flag {
	return []cli.Flag{
		ConfigFlag,
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
		&cli.StringFlag{
			Name:    "encoding",
			Aliases: []string{"e"},
			Value:   "native",
			Usage:   "Encoding to use (native, json, xml, geojson, gpx)",
		},
		&cli.BoolFlag{
			Name:  "http-tracing",
			Value: false,
			Usage: "Log all http calls (warning: no effort is made to mask log ids, keys, and other sensitive information)",
		},
		&cli.DurationFlag{
			Name:    "timeout",
			Aliases: []string{"t"},
			Value:   time.Millisecond * 10000,
			Usage:   "Timeout duration (eg, 1ms, 2s, 5m, 3h)",
		},
	}
}()

var gravlCommands = func() []*cli.Command {
	return []*cli.Command{
		analysis.Command,
		cyclinganalytics.Command,
		gnis.Command,
		gpx.Command,
		commandsCommand,
		noaa.Command,
		openweather.Command,
		rwgps.Command,
		srtm.Command,
		store.Command,
		strava.Command,
		version.Command,
		visualcrossing.Command,
		wta.Command,
	}
}()

// Run the gravl application
func Run(args []string) error {
	app := &cli.App{
		Name:     "gravl",
		HelpName: "gravl",
		Flags:    flags,
		Commands: gravlCommands,
		Before:   commands.Before(initLogging, initEncoding, initConfig),
		ExitErrHandler: func(c *cli.Context, err error) {
			if err == nil {
				return
			}
			log.Error().Err(err).Msg(c.App.Name)
		},
	}
	return app.RunContext(context.Background(), args)
}
