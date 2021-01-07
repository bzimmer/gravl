package strava

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var Command = &cli.Command{
	Name:     "strava",
	Category: "activity",
	Usage:    "Query Strava for rides and routes",
	Flags:    AuthFlags,
	Subcommands: []*cli.Command{
		activitiesCommand,
		activityCommand,
		athleteCommand,
		exportCommand,
		fitnessCommand,
		refreshCommand,
		routeCommand,
		routesCommand,
		streamsCommand,
		storeCommand,
	},
}

var AuthFlags = []cli.Flag{
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "strava.client-id",
		Usage: "API key for Strava API",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "strava.client-secret",
		Usage: "API secret for Strava API",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "strava.access-token",
		Usage: "Access token for Strava API",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "strava.refresh-token",
		Usage: "Refresh token for Strava API",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "strava.username",
		Usage: "Username for the Strava website",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "strava.password",
		Usage: "Password for the Strava website",
	}),
}
