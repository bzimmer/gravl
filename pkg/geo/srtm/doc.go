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

	"github.com/twpayne/go-geom"
)

func Example() {
	// Barlow Pass
	point := geom.NewPointFlat(geom.XY, []float64{-121.4440005, 48.0264959})
	client, err := NewClient(WithStorageLocation("/tmp"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	elevation, err := client.Elevation.Elevation(ctx, point)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(elevation)
	os.Exit(0)
}
