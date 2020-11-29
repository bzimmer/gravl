package gravl

import (
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg"
)

var versionCommand = &cli.Command{
	Name:     "version",
	Category: "api",
	Usage:    "Version",
	Action: func(c *cli.Context) error {
		return encoder.Encode(map[string]string{
			"version":    pkg.BuildVersion,
			"timestamp":  pkg.BuildTime,
			"user-agent": pkg.UserAgent,
		})
	},
}
