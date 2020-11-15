package gravl

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bzimmer/transport"
	"github.com/markbates/goth"
	auth "github.com/markbates/goth/providers/strava"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/common/route"
	"github.com/bzimmer/gravl/pkg/strava"
)

func newStravaAuthProvider(c *cli.Context, callback string) goth.Provider {
	provider := auth.New(
		c.String("strava.api-key"), c.String("strava.api-secret"), callback,
		// appears to be a bug where scope varargs do not work properly
		"read_all,profile:read_all,activity:read_all")
	t := http.DefaultTransport
	if c.Bool("http-tracing") {
		t = &transport.VerboseTransport{
			Transport: t,
		}
	}
	provider.HTTPClient = &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
	return provider
}

var stravaCommand = &cli.Command{
	Name:     "strava",
	Category: "route",
	Usage:    "Query Strava for rides and routes",
	Flags: []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "strava.api-key",
			Usage: "API key for Strava API",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "strava.api-secret",
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
		&cli.BoolFlag{
			Name:    "trip",
			Aliases: []string{"t"},
			Value:   false,
			Usage:   "Trip",
		},
		&cli.BoolFlag{
			Name:    "route",
			Aliases: []string{"r"},
			Value:   false,
			Usage:   "Route",
		},
		&cli.BoolFlag{
			Name:    "athlete",
			Aliases: []string{"a"},
			Value:   false,
			Usage:   "Athlete",
		},
		&cli.BoolFlag{
			Name:    "activities",
			Aliases: []string{"c"},
			Value:   false,
			Usage:   "Activities",
		},
		&cli.IntFlag{
			Name:    "count",
			Aliases: []string{"N"},
			Value:   10,
			Usage:   "Count",
		},
		&cli.BoolFlag{
			Name:  "refresh",
			Value: false,
			Usage: "Refresh",
		},
	},
	Action: func(c *cli.Context) error {
		client, err := strava.NewClient(
			strava.WithHTTPTracing(c.Bool("http-tracing")),
			strava.WithAPICredentials(
				c.String("strava.access-token"),
				c.String("strava.refresh-token")),
			strava.WithProvider(newStravaAuthProvider(c, "")))
		if err != nil {
			return err
		}
		r := c.Bool("route")
		a := c.Bool("activity")
		if r || a {
			args := c.Args()
			var rte *route.Route
			for i := 0; i < args.Len(); i++ {
				activityID, err := strconv.ParseInt(args.Get(i), 0, 64)
				if r {
					rte, err = client.Route.Route(c.Context, activityID)
				} else if a {
					rte, err = client.Activity.Route(c.Context, activityID)
				}
				if err != nil {
					return err
				}
				encoder.Encode(rte)
			}
		}
		if c.Bool("athlete") {
			ath, err := client.Athlete.Athlete(c.Context)
			if err != nil {
				return err
			}
			encoder.Encode(ath)
		}
		if c.Bool("refresh") {
			tokens, err := client.Auth.Refresh(c.Context)
			if err != nil {
				return err
			}
			encoder.Encode(tokens)
		}
		if c.Bool("activities") {
			acts, err := client.Activity.Activities(c.Context, c.Int("count"))
			if err != nil {
				return err
			}
			for _, act := range acts {
				encoder.Encode(act)
			}
		}
		return nil
	},
}
