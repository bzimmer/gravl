//go:generate go run ../../cmd/genreadme/genreadme.go

/*
Package gnis provides a client to the Geographic Names database provided by the USGS.

The definitions for feature classes can be found (here) https://geonames.usgs.gov/apex/f?p=GNISPQ:8::::::

To use:

	import (
		"context"
		"log"
		"os"

		"github.com/bzimmer/gravl/pkg/gnis"
	)

	func main() {
		client, err := gnis.NewClient()
		if err != nil {
			log.Fatal(err)
		}
		ctx := context.Background()
		features, err := client.GeoNames.Query(ctx, "WA")
		if err != nil {
			log.Fatal(err)
		}
		log.Println(features)
		os.Exit(0)
	}

*/
package gnis
