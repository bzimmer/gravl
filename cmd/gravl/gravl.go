package main

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/analysis/passes/ageride"
	"github.com/bzimmer/gravl/pkg/analysis/passes/benford"
	"github.com/bzimmer/gravl/pkg/analysis/passes/climbing"
	"github.com/bzimmer/gravl/pkg/analysis/passes/cluster"
	"github.com/bzimmer/gravl/pkg/analysis/passes/eddington"
	"github.com/bzimmer/gravl/pkg/analysis/passes/festive500"
	"github.com/bzimmer/gravl/pkg/analysis/passes/forecast"
	"github.com/bzimmer/gravl/pkg/analysis/passes/hourrecord"
	"github.com/bzimmer/gravl/pkg/analysis/passes/koms"
	"github.com/bzimmer/gravl/pkg/analysis/passes/pythagorean"
	"github.com/bzimmer/gravl/pkg/analysis/passes/rolling"
	"github.com/bzimmer/gravl/pkg/analysis/passes/splat"
	"github.com/bzimmer/gravl/pkg/analysis/passes/staticmap"
	"github.com/bzimmer/gravl/pkg/analysis/passes/totals"
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
	"github.com/bzimmer/gravl/pkg/commands/manual"
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
	initAnalysis := func(c *cli.Context) error {
		analysis.Add(ageride.New(), false)
		analysis.Add(benford.New(), false)
		analysis.Add(climbing.New(), true)
		analysis.Add(cluster.New(), false)
		analysis.Add(eddington.New(), true)
		analysis.Add(festive500.New(), true)
		analysis.Add(forecast.New(), false)
		analysis.Add(hourrecord.New(), true)
		analysis.Add(koms.New(), true)
		analysis.Add(pythagorean.New(), true)
		analysis.Add(rolling.New(), true)
		analysis.Add(splat.New(), false)
		analysis.Add(staticmap.New(), false)
		analysis.Add(totals.New(), true)
		return nil
	}
	commands := []*cli.Command{
		analysis.Command,
		cyclinganalytics.Command,
		gnis.Command,
		gpx.Command,
		manual.Commands,
		manual.Manual,
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
		Name:        "gravl",
		HelpName:    "gravl",
		Usage:       "CLI for activity related analysis, exploration, & planning",
		Description: "Activity related analysis, exploration, & planning",
		Flags:       gravl.Flags("gravl.yaml"),
		Commands:    commands,
		Before:      gravl.Befores(gravl.InitLogging(), initEncoding, gravl.InitConfig(), initAnalysis),
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
