package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/commands/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/pkg/commands/activity/rwgps"
	"github.com/bzimmer/gravl/pkg/commands/activity/strava"
	"github.com/bzimmer/gravl/pkg/commands/activity/wta"
	"github.com/bzimmer/gravl/pkg/commands/analysis"
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
	commands := []*cli.Command{
		analysis.Command,
		cyclinganalytics.Command,
		gnis.Command,
		gpx.Command,
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
	app := gravl.App("gravl", commands)
	if err := app.RunContext(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
