package visualcrossing_test

import (
	"time"

	"github.com/bzimmer/gravl/pkg/wx/visualcrossing"
	"github.com/bzimmer/httpwares"
)

func newClient(status int, filename string) (*visualcrossing.Client, error) {
	return visualcrossing.NewClient(
		visualcrossing.WithTransport(&httpwares.TestDataTransport{
			Status:      status,
			Filename:    filename,
			ContentType: "application/json",
		}),
		visualcrossing.WithTokenCredentials("fooKey", "", time.Time{}))
}
