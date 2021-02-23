package main

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/commands/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/pkg/commands/activity/rwgps"
	"github.com/bzimmer/gravl/pkg/commands/activity/strava"
	"github.com/bzimmer/gravl/pkg/commands/activity/wta"
	"github.com/bzimmer/gravl/pkg/commands/activity/zwift"
	"github.com/bzimmer/gravl/pkg/commands/analysis"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/commands/geo/gnis"
	"github.com/bzimmer/gravl/pkg/commands/geo/gpx"
	"github.com/bzimmer/gravl/pkg/commands/geo/srtm"
	"github.com/bzimmer/gravl/pkg/commands/gravl"
	"github.com/bzimmer/gravl/pkg/commands/store"
	"github.com/bzimmer/gravl/pkg/commands/version"
	"github.com/bzimmer/gravl/pkg/commands/wx/noaa"
	"github.com/bzimmer/gravl/pkg/commands/wx/openweather"
	"github.com/bzimmer/gravl/pkg/commands/wx/visualcrossing"
)

func main() {
	initEncoding := gravl.InitEncoding(
		func(c *cli.Context) encoding.Encoder {
			return encoding.GPX(c.App.Writer, c.Bool("compact"))
		},
		func(c *cli.Context) encoding.Encoder {
			return encoding.GeoJSON(c.App.Writer, c.Bool("compact"))
		},
		func(c *cli.Context) encoding.Encoder {
			return encoding.Named(c.App.Writer, c.Bool("compact"))
		},
	)
	commands := []*cli.Command{
		analysis.Command,
		cyclinganalytics.Command,
		gnis.Command,
		gpx.Command,
		gravl.Commands,
		noaa.Command,
		openweather.Command,
		rwgps.Command,
		srtm.Command,
		store.Command,
		strava.Command,
		version.Command,
		visualcrossing.Command,
		wta.Command,
		zwift.Command,
	}
	app := &cli.App{
		Name:     "gravl",
		HelpName: "gravl",
		Usage:    "Clients for activty-related services and an extensible analysis framework for activities",
		Flags:    gravl.Flags("gravl.yaml"),
		Commands: commands,
		Before:   gravl.Befores(gravl.InitLogging(), initEncoding, gravl.InitConfig()),
		ExitErrHandler: func(c *cli.Context, err error) {
			if err == nil {
				return
			}
			log.Error().Err(err).Msg(c.App.Name)
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := app.RunContext(ctx, os.Args); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
