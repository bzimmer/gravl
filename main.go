package main

import (
	"os"

	"github.com/bzimmer/gravl/cmd/gravl"
)

func main() {
	err := gravl.Run()
	if err != nil {
		os.Exit(1)
	}
}
