package gravl

import (
	"context"
	stdlog "log"
	"os"
	"path"
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
	"github.com/bzimmer/gravl/pkg/commands/geo/srtm"
	"github.com/bzimmer/gravl/pkg/commands/serve"
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

func initConfig(c *cli.Context) error {
	cfg := c.String("config")
	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		log.Warn().
			Str("path", cfg).
			Msg("unable to find config file")
		return nil
	}
	config := func() (altsrc.InputSourceContext, error) {
		return altsrc.NewYamlSourceFromFile(cfg)
	}
	for _, cmd := range c.App.Commands {
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

var flags = func() []cli.Flag {
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
		&cli.StringFlag{
			Name:    "encoding",
			Aliases: []string{"e"},
			Value:   "native",
			Usage:   "Encoding to use (native, json, xml, geojson, gpx)",
		},
		&cli.BoolFlag{
			Name:  "http-tracing",
			Value: false,
			Usage: "Log all http calls (warning: this will log ids, keys, and other sensitive information)",
		},
		&cli.PathFlag{
			Name:      "config",
			Value:     config,
			TakesFile: true,
			Usage:     "File containing configuration settings",
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
		cyclinganalytics.Command,
		gnis.Command,
		noaa.Command,
		openweather.Command,
		analysis.Command,
		rwgps.Command,
		serve.Command,
		srtm.Command,
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
			log.Error().Err(err).Msg(c.App.Name)
		},
	}
	return app.RunContext(context.Background(), args)
}
