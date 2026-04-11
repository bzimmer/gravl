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
	// BuildBuilder of the package
	BuildBuilder = "local"
	// BuildCommit of the package
	BuildCommit = "unknown"
	// UserAgent of the package
	UserAgent = fmt.Sprintf("gravl/%s (https://github.com/bzimmer/gravl)", BuildVersion)
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Version",
		Action: func(c *cli.Context) error {
			m := map[string]string{
				"builder":    BuildBuilder,
				"commit":     BuildCommit,
				"timestamp":  BuildTime,
				"user-agent": UserAgent,
				"version":    BuildVersion,
			}
			if info, ok := debug.ReadBuildInfo(); ok {
				m["go"] = info.GoVersion
			}
			return gravl.Runtime(c).Encoder.Encode(m)
		},
	}
}
