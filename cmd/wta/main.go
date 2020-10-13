package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/bzimmer/wta"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func run(ctx *cli.Context) error {
	c := colly.NewCollector(
		colly.AllowedDomains("wta.org", "www.wta.org"),
	)

	if ctx.IsSet("serve") {
		r := gin.New()

		r.Use(logger.SetLogger())
		r.Use(gin.Recovery())

		r.GET("/regions/", wta.RegionsHandler())
		r.GET("/reports/", wta.TripReportsHandler(c))
		r.GET("/reports/:reporter", wta.TripReportsHandler(c))

		address := fmt.Sprintf("0.0.0.0:%d", ctx.Value("port").(int))
		if err := r.Run(address); err != nil {
			return err
		}
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	args := ctx.Args()
	for i := 0; i < args.Len(); i++ {
		q, err := wta.Query(args.Get(i))
		if err != nil {
			return err
		}

		reports, err := wta.GetTripReports(c, q.String())
		if err != nil {
			return err
		}

		for _, tr := range reports {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%d\t%s\t%s\t%s", tr.Title, tr.Votes, tr.HikeDate.Format("2006-01-02"), tr.Region, tr.Report))
		}
	}

	return nil
}

func main() {
	app := &cli.App{
		Name:   "blexp",
		Usage:  "submit expenses from the cli",
		Action: run,
		Flags: []cli.Flag{
			// &cli.BoolFlag{
			// 	Name:    "json",
			// 	Value:   false,
			// 	Aliases: []string{"j"},
			// 	Usage:   "JSON output",
			// },
			&cli.BoolFlag{
				Name:    "tab",
				Value:   true,
				Aliases: []string{"t"},
				Usage:   "Tabular output",
			},
			&cli.BoolFlag{
				Name:    "serve",
				Value:   false,
				Aliases: []string{"s"},
				Usage:   "Serve the results via a REST API listening on `PORT`",
			},
			&cli.IntFlag{
				Name:    "port",
				Value:   1122,
				Aliases: []string{"p"},
				Usage:   "Port",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	os.Exit(0)
}
