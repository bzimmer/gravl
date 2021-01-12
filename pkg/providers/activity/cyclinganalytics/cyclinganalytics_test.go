package cyclinganalytics_test

import (
	"time"

	"github.com/bzimmer/gravl/pkg/providers/activity/cyclinganalytics"
	"github.com/bzimmer/httpwares"
)

func newClient(status int, filename string) (*cyclinganalytics.Client, error) {
	return cyclinganalytics.NewClient(
		cyclinganalytics.WithTransport(&httpwares.TestDataTransport{
			Status:      status,
			Filename:    filename,
			ContentType: "application/json"}),
		cyclinganalytics.WithHTTPTracing(false),
		cyclinganalytics.WithTokenCredentials("fooKey", "barToken", time.Time{}))
}
