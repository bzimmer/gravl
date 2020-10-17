package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	gn "github.com/bzimmer/wta/pkg/gnis"
	"github.com/bzimmer/wta/pkg/wta"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

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
	r := wta.NewRouter(wta.NewCollector())

	port := ctx.Value("port").(int)
	address := fmt.Sprintf("0.0.0.0:%d", port)
	if err := r.Run(address); err != nil {
		return err
	}
	return nil
}

func gnis(ctx *cli.Context) error {
	g := gn.New()
	encoder := json.NewEncoder(os.Stdout)
	if !ctx.IsSet("compact") {
		encoder.SetIndent("", " ")
	}
	encoder.SetEscapeHTML(false)
	for _, arg := range ctx.Args().Slice() {
		log.Info().Str("filename", arg)
		features, err := g.ParseFile(arg)
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

func list(ctx *cli.Context) error {
	args := ctx.Args().Slice()
	if len(args) == 0 {
		// query the most recent if no reporter specified
		args = append(args, "")
	}

	c := wta.NewCollector()
	reports := make([]wta.TripReport, 0)
	for _, arg := range args {
		q := wta.Query(arg)
		tr, err := wta.GetTripReports(c, q.String())
		if err != nil {
			return err
		}
		for _, r := range tr {
			reports = append(reports, r)
		}
	}

	encoder := json.NewEncoder(os.Stdout)
	if !ctx.IsSet("compact") {
		encoder.SetIndent("", " ")
	}
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(reports)
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
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	os.Exit(0)
}
