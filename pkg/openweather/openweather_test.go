package openweather_test

import (
	ow "github.com/bzimmer/gravl/pkg/openweather"
	"github.com/bzimmer/httpwares"
)

func newClient(status int, filename string) (*ow.Client, error) {
	return ow.NewClient(
		ow.WithTransport(&httpwares.TestDataTransport{
			Status:      status,
			Filename:    filename,
			ContentType: "application/json",
		}),
		ow.WithAPIKey("fooKey"))
}
