package gpx

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/twpayne/go-gpx"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/geo"
)

func info(c *cli.Context) error {
	for _, arg := range c.Args().Slice() {
		fp, err := os.Open(arg)
		if err != nil {
			return err
		}
		defer func() {
			if err = fp.Close(); err != nil {
				log.Error().Err(err).Msg("info")
			}
		}()
		x, err := gpx.Read(fp)
		if err != nil {
			return err
		}
		s := geo.SummarizeTracks(x)
		if s.Tracks > 0 {
			s.Filename = arg
			if err := encoding.Encode(s); err != nil {
				return err
			}
		}
		s = geo.SummarizeRoutes(x)
		if s.Routes > 0 {
			s.Filename = arg
			if err := encoding.Encode(s); err != nil {
				return err
			}
		}
	}
	return nil
}

var infoCommand = &cli.Command{
	Name:      "info",
	Usage:     "Return basic statistics about a GPX file",
	ArgsUsage: "GPX_FILE (...)",
	Action:    info,
}

var Command = &cli.Command{
	Name:        "gpx",
	Usage:       "gpx",
	Category:    "geo",
	Subcommands: []*cli.Command{infoCommand},
}
