package noaa_test

import (
	"github.com/bzimmer/gravl/pkg/noaa"
	"github.com/bzimmer/transport"
)

func newClient(status int, filename string) (*noaa.Client, error) {
	return noaa.NewClient(
		noaa.WithTransport(&transport.TestDataTransport{
			Status:      status,
			Filename:    filename,
			ContentType: "application/json",
		}),
	)
}
