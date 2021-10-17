package version

import (
	"fmt"

	"github.com/bzimmer/gravl/pkg"
	"github.com/urfave/cli/v2"
)

var (
	// BuildVersion of the package
	BuildVersion = "development"
	// BuildTime of the package
	BuildTime = "now"
	// UserAgent of the package
	UserAgent = fmt.Sprintf("gravl/%s (https://github.com/bzimmer/gravl)", BuildVersion)
	// PackageName is the name of the package
	PackageName = "com.github.bzimmer.gravl"
)

var Command = &cli.Command{
	Name:  "version",
	Usage: "Version",
	Action: func(c *cli.Context) error {
		return pkg.Runtime(c).Encoder.Encode(map[string]string{
			"version":    BuildVersion,
			"timestamp":  BuildTime,
			"user-agent": UserAgent,
			"config":     c.String("config"),
		})
	},
}
