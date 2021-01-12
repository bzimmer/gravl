// Package gnis provides a client to the Geographic Names database provided by the USGS.

// The definitions for feature classes can be found (here) https://geonames.usgs.gov/apex/f?p=GNISPQ:8::::::
package gnis

import (
	"context"
	"log"
	"os"
)

func Example() {
	client, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	features, err := client.GeoNames.Query(ctx, "WA")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(features[0])
	os.Exit(0)
}
