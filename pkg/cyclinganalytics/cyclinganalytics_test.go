package cyclinganalytics_test

import (
	"time"

	"github.com/bzimmer/gravl/pkg/cyclinganalytics"
	"github.com/bzimmer/httpwares"
)

func newClient(status int, filename string) (*cyclinganalytics.Client, error) {
	return cyclinganalytics.NewClient(
		cyclinganalytics.WithTransport(&httpwares.TestDataTransport{
			Status:      status,
			Filename:    filename,
			ContentType: "application/json"}),
		cyclinganalytics.WithHTTPTracing(true),
		cyclinganalytics.WithTokenCredentials("fooKey", "barToken", time.Time{}))
}
