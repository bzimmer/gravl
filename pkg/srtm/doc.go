/*
Package srtm provides a client to access the SRTM elevation dataset.

Information about the Shuttle Radar Topology Mission can be found at NASA's
(site) https://www2.jpl.nasa.gov/srtm/. This client is itself a wrapper around
the (go-elevations) https://github.com/tkrajina/go-elevations/ client library.
*/
package srtm

import (
	"context"
	"log"
	"os"
)

func Example() {
	// Barlow Pass
	longitude, latitude := -121.4440005, 48.0264959
	client, err := NewClient(
		WithStorageLocation("/path/to/storage/directory"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	elevation, err := client.Elevation.Elevation(ctx, longitude, latitude)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(elevation)
	os.Exit(0)
}
