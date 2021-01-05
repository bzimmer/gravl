// Package openweather provides a client to access the OpenWeather's API.
package openweather

import (
	"context"
	"log"
	"time"

	"github.com/twpayne/go-geom"
)

func Example() {
	client, err := NewClient(
		WithTokenCredentials("openweather.access-token", "", time.Time{}))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	fcst, err := client.Forecast.Forecast(ctx,
		ForecastOptions{
			Units: UnitsMetric,
			Point: geom.NewPointFlat(geom.XY, []float64{-122.2992, 48.82})})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fcst)
}
