package openweather_test

import (
	"time"

	"github.com/bzimmer/gravl/pkg/openweather"
	"github.com/bzimmer/httpwares"
)

func newClient(status int, filename string) (*openweather.Client, error) {
	return openweather.NewClient(
		openweather.WithTransport(&httpwares.TestDataTransport{
			Status:      status,
			Filename:    filename,
			ContentType: "application/json",
		}),
		openweather.WithTokenCredentials("fooKey", "", time.Time{}))
}
