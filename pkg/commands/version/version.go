package version

import (
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
)

var Command = &cli.Command{
	Name:  "version",
	Usage: "Version",
	Action: func(c *cli.Context) error {
		return encoding.For(c).Encode(map[string]string{
			"version":    pkg.BuildVersion,
			"timestamp":  pkg.BuildTime,
			"user-agent": pkg.UserAgent,
			"config":     c.String("config"),
		})
	},
}
