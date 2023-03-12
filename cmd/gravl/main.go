package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/armon/go-metrics"
	"github.com/bzimmer/activity"
	"github.com/bzimmer/manual"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/activity/qp"
	"github.com/bzimmer/gravl/activity/rwgps"
	"github.com/bzimmer/gravl/activity/strava"
	"github.com/bzimmer/gravl/activity/zwift"
	"github.com/bzimmer/gravl/eval/antonmedv"
	"github.com/bzimmer/gravl/version"
)

func initSignal(cancel context.CancelFunc) cli.BeforeFunc {
	return func(c *cli.Context) error {
		go func() {
			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, os.Interrupt)
			select {
			case <-sigc:
				log.Info().Msg("canceling...")
				cancel()
			case <-c.Context.Done():
			}
		}()
		return nil
	}
}

func initQP(c *cli.Context) error {
	// strava
	gravl.Runtime(c).Exporters[strava.Provider] = func(c *cli.Context) (activity.Exporter, error) {
		if err := strava.Before(c); err != nil {
			return nil, err
		}
		return gravl.Runtime(c).Strava.Exporter(), nil
	}
	gravl.Runtime(c).Uploaders[strava.Provider] = func(c *cli.Context) (activity.Uploader, error) {
		if err := strava.Before(c); err != nil {
			return nil, err
		}
		return gravl.Runtime(c).Strava.Uploader(), nil
	}
	// cyclinganalytics
	gravl.Runtime(c).Uploaders[cyclinganalytics.Provider] = func(c *cli.Context) (activity.Uploader, error) {
		if err := cyclinganalytics.Before(c); err != nil {
			return nil, err
		}
		return gravl.Runtime(c).CyclingAnalytics.Uploader(), nil
	}
	// zwift
	gravl.Runtime(c).Exporters[zwift.Provider] = func(c *cli.Context) (activity.Exporter, error) {
		if err := zwift.Before(c); err != nil {
			return nil, err
		}
		return gravl.Runtime(c).Zwift.Exporter(), nil
	}
	return nil
}

func initRuntime(c *cli.Context) error {
	writer := io.Discard
	if c.Bool("json") {
		writer = c.App.Writer
	}

	cfg := metrics.DefaultConfig(c.App.Name)
	cfg.EnableRuntimeMetrics = false
	cfg.TimerGranularity = time.Second
	sink := metrics.NewInmemSink(time.Hour*24, time.Hour*24)
	metric, err := metrics.New(cfg, sink)
	if err != nil {
		return err
	}

	c.App.Metadata[gravl.RuntimeKey] = &gravl.Rt{
		Start:     time.Now(),
		Encoder:   json.NewEncoder(writer),
		Filterer:  antonmedv.Filterer,
		Evaluator: antonmedv.Evaluator,
		Sink:      sink,
		Metrics:   metric,
		Fs:        afero.NewOsFs(),
		Uploaders: make(map[string]gravl.UploaderFunc),
		Exporters: make(map[string]gravl.ExporterFunc),
		Endpoints: make(map[string]oauth2.Endpoint),
	}
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
	return nil
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
			Name:     "json",
			Aliases:  []string{"j"},
			Usage:    "Emit all results as JSON and print to stdout",
			Value:    false,
			Required: false,
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
		manual.Manual(),
		manual.Commands(),
		manual.EnvVars(),
		qp.Command(),
		rwgps.Command(),
		strava.Command(),
		version.Command(),
		zwift.Command(),
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	app := &cli.App{
		Name:        "gravl",
		HelpName:    "gravl",
		Usage:       "command line access to activity platforms",
		Description: "command line access to activity platforms",
		Flags:       flags(),
		Commands:    commands(),
		Before:      gravl.Befores(initSignal(cancel), initLogging, initRuntime, initQP),
		After: func(c *cli.Context) error {
			t := gravl.Runtime(c).Start
			met := gravl.Runtime(c).Metrics
			met.AddSample([]string{"runtime"}, float32(time.Since(t).Seconds()))
			return gravl.Stats(c)
		},
	}
	var err error
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				log.Error().Err(v).Msg(app.Name)
			case string:
				log.Error().Err(errors.New(v)).Msg(app.Name)
			default:
				log.Error().Err(fmt.Errorf("%v", v)).Msg(app.Name)
			}
			os.Exit(1)
		}
		if err != nil {
			log.Error().Err(err).Msg(app.Name)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	err = app.RunContext(ctx, os.Args)
}
