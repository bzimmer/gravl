package openweather_test

import (
	"time"

	"github.com/bzimmer/httpwares"

	"github.com/bzimmer/gravl/pkg/providers/wx/openweather"
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
