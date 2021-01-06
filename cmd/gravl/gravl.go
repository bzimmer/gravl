package main

import (
	"os"

	"github.com/bzimmer/gravl/pkg/commands/gravl"
)

func main() {
	if err := gravl.Run(os.Args); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
