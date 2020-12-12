//go:generate go run ../../cmd/genreadme/genreadme.go

/*
Package openweather provides a client to access the OpenWeather's API.

An example:

	import (
		"context"
		"log"
		"time"

		"github.com/bzimmer/gravl/pkg/openweather"
	)

	func main() {
		client, err := openweather.NewClient(
			WithTokenCredentials("openweather.access-token", "", time.Time{}))
		if err != nil {
			log.Fatal(err)
		}
		ctx := context.Background()
		fcst, err := client.Forecast.Forecast(ctx,
			openweather.ForecastOptions{
				Units: openweather.UnitsMetric,
				Point: geom.NewPointFlat(geom.XY, []float64{longitude, latitude})})
		if err != nil {
			log.Fatal(err)
		}
		log.Println(fcst)
	}
*/
package openweather
