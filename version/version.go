package version

import (
	"fmt"
	"runtime/debug"

	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl"
)

var (
	// BuildVersion of the package
	BuildVersion = "development"
	// BuildTime of the package
	BuildTime = "now"
	// UserAgent of the package
	UserAgent = fmt.Sprintf("gravl/%s (https://github.com/bzimmer/gravl)", BuildVersion)
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Version",
		Action: func(c *cli.Context) error {
			m := map[string]string{
				"version":    BuildVersion,
				"timestamp":  BuildTime,
				"user-agent": UserAgent,
			}
			if info, ok := debug.ReadBuildInfo(); ok {
				m["go"] = info.GoVersion
			}
			return gravl.Runtime(c).Encoder.Encode(m)
		},
	}
}
