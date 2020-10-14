package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

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
			NoColor: !ctx.IsSet("color"),
		},
	)
	log.Info().Msg("configured logging")
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

	if ctx.IsSet("tab") {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer w.Flush()
		for _, tr := range reports {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%d\t%s\t%s\t%s", tr.Title, tr.Votes, tr.HikeDate.Format("2006-01-02"), tr.Region, tr.Report))
		}
		return nil
	}

	// json is the default
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", " ")
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
				Name:  "color",
				Value: false,
				Usage: "Colorize the console log",
			},
		},
		Commands: []*cli.Command{
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
				Usage:   "List the results to the console (in JSON by default)",
				Action:  list,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "tab",
						Value:   false,
						Aliases: []string{"t"},
						Usage:   "List results in tabular output",
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
