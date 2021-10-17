package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/armon/go-metrics"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/pkg/activity/qp"
	"github.com/bzimmer/gravl/pkg/activity/rwgps"
	"github.com/bzimmer/gravl/pkg/activity/strava"
	"github.com/bzimmer/gravl/pkg/activity/zwift"
	"github.com/bzimmer/gravl/pkg/eval/antonmedv"
	"github.com/bzimmer/gravl/pkg/manual"
	"github.com/bzimmer/gravl/pkg/version"
)

func initRuntime() cli.BeforeFunc {
	return func(c *cli.Context) error {
		var enc pkg.Encoder
		compact := c.Bool("compact")
		switch c.String("encoding") {
		case "spew":
			enc = pkg.Spew(c.App.Writer)
		// case "geojson":
		// 	enc = gravl.GeoJSON(c.App.Writer, compact)
		case "xml":
			enc = pkg.XML(c.App.Writer, compact)
		case "json":
			enc = pkg.JSON(c.App.Writer, compact)
		default:
			enc = pkg.Blackhole()
		}

		cfg := metrics.DefaultConfig("gravl")
		cfg.EnableRuntimeMetrics = false
		cfg.TimerGranularity = time.Second
		sink := metrics.NewInmemSink(time.Hour*24, time.Hour*24)
		metric, err := metrics.New(cfg, sink)
		if err != nil {
			return err
		}

		c.App.Metadata[pkg.RuntimeKey] = &pkg.Rt{
			Start:     time.Now(),
			Encoder:   enc,
			Mapper:    antonmedv.Mapper,
			Filterer:  antonmedv.Filterer,
			Evaluator: antonmedv.Evaluator,
			Sink:      sink,
			Metrics:   metric,
		}
		return nil
	}
}

func initLogging() cli.BeforeFunc {
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
		return nil
	}
}

func flags() []cli.Flag {
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
			Value:   "",
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
	}
}

func commands() []*cli.Command {
	return []*cli.Command{
		cyclinganalytics.Command(),
		manual.Command(),
		manual.Commands(),
		qp.Command(),
		rwgps.Command(),
		strava.Command(),
		version.Command(),
		zwift.Command(),
	}
}

func run() error {
	app := &cli.App{
		Name:        "gravl",
		HelpName:    "gravl",
		Usage:       "CLI for activity related analysis, exploration, & planning",
		Description: "Activity related analysis, exploration, & planning",
		Flags:       flags(),
		Commands:    commands(),
		After: func(c *cli.Context) error {
			t := pkg.Runtime(c).Start
			met := pkg.Runtime(c).Metrics
			met.AddSample([]string{"runtime"}, float32(time.Since(t).Seconds()))
			return pkg.Stats(c)
		},
		Before: pkg.Befores(initLogging(), initRuntime()),
		ExitErrHandler: func(c *cli.Context, err error) {
			if err == nil {
				return
			}
			log.Error().Err(err).Msg(c.App.Name)
		},
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			log.Info().Msg("canceling...")
			cancel()
		case <-ctx.Done():
		}
		<-c
		os.Exit(2)
	}()
	return app.RunContext(ctx, os.Args)
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
