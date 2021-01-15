package gpx

import (
	"os"

	"github.com/urfave/cli/v2"

	ggpx "github.com/twpayne/go-gpx"

	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/geo"
)

func info(c *cli.Context) error {
	n := c.NArg()
	for i := 0; i < n; i++ {
		arg := c.Args().Get(i)
		fp, err := os.Open(arg)
		if err != nil {
			return err
		}
		defer fp.Close()
		x, err := ggpx.Read(fp)
		if err != nil {
			return err
		}
		s := geo.SummarizeTracks(x)
		if s.Tracks > 0 {
			if err := encoding.Encode(s); err != nil {
				return err
			}
		}
		s = geo.SummarizeRoutes(x)
		if s.Routes > 0 {
			if err := encoding.Encode(s); err != nil {
				return err
			}
		}
	}
	return nil
}

var infoCommand = &cli.Command{
	Name:   "info",
	Usage:  "Return basic statistics about a GPX file",
	Action: info,
}

var Command = &cli.Command{
	Name:        "gpx",
	Usage:       "gpx",
	Category:    "geo",
	Subcommands: []*cli.Command{infoCommand},
}
