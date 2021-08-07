package gravl

import (
	"errors"
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
	"github.com/bzimmer/gravl/pkg/commands/encoding"
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

// Befores combines multiple `cli.BeforeFunc`s into a single `cli.BeforeFunc`
func Befores(bfs ...cli.BeforeFunc) cli.BeforeFunc {
	return func(c *cli.Context) error {
		for _, fn := range bfs {
			if fn == nil {
				continue
			}
			if err := fn(c); err != nil {
				return err
			}
		}
		return nil
	}
}

func InitConfig() cli.BeforeFunc {
	return func(c *cli.Context) error {
		cfg := c.String("config")
		if _, err := os.Stat(cfg); os.IsNotExist(err) {
			log.Error().Str("path", cfg).Msg("unable to find config file")
			return errors.New("invalid config file")
		}
		config := func() (altsrc.InputSourceContext, error) {
			return altsrc.NewYamlSourceFromFile(cfg)
		}
		// configure the application flags
		if err := altsrc.InitInputSource(c.App.Flags, config)(c); err != nil {
			return err
		}
		// configure the subcommand flags
		for _, cmd := range flatten(c.App.Commands) {
			cmd.Before = Befores(altsrc.InitInputSource(cmd.Flags, config), cmd.Before)
		}
		return nil
	}
}

func InitEncoding() cli.BeforeFunc {
	return func(c *cli.Context) error {
		var enc encoding.Encoder
		compact := c.Bool("compact")
		switch c.String("encoding") {
		case "spew":
			enc = encoding.Spew(c.App.Writer)
		case "geojson":
			enc = encoding.GeoJSON(c.App.Writer, compact)
		case "xml":
			enc = encoding.XML(c.App.Writer, compact)
		case "named":
			enc = encoding.Named(c.App.Writer, compact)
		default:
			enc = encoding.JSON(c.App.Writer, compact)
		}
		c.App.Metadata["enc"] = enc
		return nil
	}
}

func InitLogging() cli.BeforeFunc {
	return func(c *cli.Context) error {
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
}

func Flags(filename string) []cli.Flag {
	config := path.Join(xdg.ConfigHome, pkg.PackageName, filename)
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
			Value:   "json",
			Usage:   "Output encoding (eg: json, xml, geojson, gpx, spew)",
		},
		&cli.BoolFlag{
			Name:  "http-tracing",
			Value: false,
			Usage: "Log all http calls (warning: no effort is made to mask log ids, keys, and other sensitive information)",
		},
		&cli.DurationFlag{
			Name:    "timeout",
			Aliases: []string{"t"},
			Value:   time.Second * 10,
			Usage:   "Timeout duration (eg, 1ms, 2s, 5m, 3h)",
		},
		&cli.StringFlag{
			Name:  "config",
			Value: config,
			Usage: "File containing configuration values",
		},
	}
}
