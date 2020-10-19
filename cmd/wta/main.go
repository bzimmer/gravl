package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	gn "github.com/bzimmer/wta/pkg/gnis"
	na "github.com/bzimmer/wta/pkg/noaa"
	"github.com/bzimmer/wta/pkg/wta"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func newEncoder(compact bool) *json.Encoder {
	encoder := json.NewEncoder(os.Stdout)
	if !compact {
		encoder.SetIndent("", " ")
	}
	encoder.SetEscapeHTML(false)
	return encoder
}

func initLogging(ctx *cli.Context) error {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.DurationFieldInteger = true
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: ctx.IsSet("monochrome"),
		},
	)
	log.Info().
		Str("url", "https://www.wta.org/").
		Str("build_version", wta.BuildVersion).
		Msg("Please support the WTA")
	return nil
}

func serve(ctx *cli.Context) error {
	log.Info().Msg("configuring to serve")
	c, err := wta.NewClient()
	if err != nil {
		return err
	}
	r := wta.NewRouter(c)

	port := ctx.Value("port").(int)
	address := fmt.Sprintf("0.0.0.0:%d", port)
	if err := r.Run(address); err != nil {
		return err
	}
	return nil
}

func gnis(ctx *cli.Context) error {
	g, err := gn.NewClient()
	if err != nil {
		return err
	}
	b := context.Background()
	encoder := newEncoder(ctx.IsSet("compact"))
	for _, arg := range ctx.Args().Slice() {
		log.Info().Str("filename", arg)
		features, err := g.GeoNames.Query(b, arg)
		if err != nil {
			return err
		}
		err = encoder.Encode(features)
		if err != nil {
			return err
		}
	}
	return nil
}

func noaa(ctx *cli.Context) error {
	c, err := na.NewClient(
		na.WithTimeout(10 * time.Second),
	)
	if err != nil {
		return err
	}
	b := context.Background()

	var (
		point    *na.GridPoint
		forecast *na.Forecast
	)
	args := ctx.Args().Slice()
	switch len(args) {
	case 2:
		lat, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			log.Error().Err(err).Send()
		}
		lng, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			log.Error().Err(err).Send()
		}
		if ctx.IsSet("point") {
			point, err = c.Points.GridPoint(b, lat, lng)
		}
		forecast, err = c.Points.Forecast(b, lat, lng)
	case 3:
		// check for -p and err if true
		wfo := args[0]
		x, _ := strconv.Atoi(args[1])
		y, _ := strconv.Atoi(args[2])
		forecast, err = c.GridPoints.Forecast(b, wfo, x, y)
	default:
		return fmt.Errorf("only 2 or 3 arguments allowed")
	}
	if err != nil {
		return err
	}
	encoder := newEncoder(ctx.IsSet("compact"))
	if ctx.IsSet("point") {
		err = encoder.Encode(point)
	} else {
		err = encoder.Encode(forecast)
	}
	if err != nil {
		return err
	}
	return nil
}

func list(ctx *cli.Context) error {
	args := ctx.Args().Slice()
	if len(args) == 0 {
		// query the most recent if no reporter specified
		args = append(args, "")
	}

	c, err := wta.NewClient()
	if err != nil {
		return err
	}
	b := context.Background()
	reports := make([]*wta.TripReport, 0)
	for _, arg := range args {
		tr, err := c.Reports.TripReports(b, arg)
		if err != nil {
			return err
		}
		for _, r := range tr {
			reports = append(reports, r)
		}
	}

	encoder := newEncoder(ctx.IsSet("compact"))
	err = encoder.Encode(reports)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	app := &cli.App{
		Name:   "wta",
		Usage:  "programmatic access to WTA trip reports",
		Before: initLogging,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "monochrome",
				Value:   false,
				Aliases: []string{"m"},
				Usage:   "Monochrome console logs",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Display the version",
				Action: func(c *cli.Context) error {
					// initLogging takes care displaying version information
					os.Exit(0)
					return nil
				},
			},
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "Serve the results via a REST API",
				Action:  serve,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "port",
						Value:   1122,
						Aliases: []string{"p"},
						Usage:   "Port on which to listen",
						EnvVars: []string{"WTA_PORT"},
					},
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List the results to the console in JSON format",
				Action:  list,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "compact",
						Value:   false,
						Aliases: []string{"c"},
						Usage:   "Compact instead of pretty-printed output",
					},
				},
			},
			{
				Name:    "gnis",
				Aliases: []string{"g"},
				Usage:   "Display GNIS in JSON format",
				Action:  gnis,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "compact",
						Value:   false,
						Aliases: []string{"c"},
						Usage:   "Compact instead of pretty-printed output",
					},
				},
			},
			{
				Name:    "noaa",
				Aliases: []string{"n"},
				Usage:   "Display NOAA forecast data",
				Action:  noaa,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "compact",
						Value:   false,
						Aliases: []string{"c"},
						Usage:   "Compact instead of pretty-printed output",
					},
					&cli.BoolFlag{
						Name:    "point",
						Value:   false,
						Aliases: []string{"p"},
						Usage:   "Return the grid point details for the coordinates",
					},
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	os.Exit(0)
}
