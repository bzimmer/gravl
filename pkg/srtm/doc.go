//go:generate go run ../../cmd/genreadme/genreadme.go

/*
Package srtm provides a client to access the SRTM elevation dataset.

Information about the Shuttle Radar Topology Mission can be found at NASA's
(site) https://www2.jpl.nasa.gov/srtm/. This client is itself a wrapper around
the (go-elevations) https://github.com/tkrajina/go-elevations/ client library.


To use:

	import (
		"context"
		"fmt"
		"os"

		"github.com/bzimmer/gravl/pkg/srtm"
	)

	func main() {
		// Barlow Pass
		longitude, latitude := -121.4440005, 48.0264959
		client, err := srtm.NewClient(
			srtm.WithStorageLocation("/path/to/storage/directory"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ctx := context.Background()
		elevation, err := client.Elevation.Elevation(ctx, longitude, latitude)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(elevation)
		os.Exit(0)
	}
*/
package srtm
