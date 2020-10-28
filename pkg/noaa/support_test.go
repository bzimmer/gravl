package noaa_test

import (
	"github.com/bzimmer/wta/pkg/common"
	"github.com/bzimmer/wta/pkg/noaa"
)

func newClient(status int, filename string) (*noaa.Client, error) {
	return noaa.NewClient(
		noaa.WithTransport(&common.TestDataTransport{
			Status:      status,
			Filename:    filename,
			ContentType: "application/json",
		}),
	)
}
