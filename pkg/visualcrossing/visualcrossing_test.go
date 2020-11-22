package visualcrossing_test

import (
	vc "github.com/bzimmer/gravl/pkg/visualcrossing"
	"github.com/bzimmer/httpwares"
)

func newClient(status int, filename string) (*vc.Client, error) {
	return vc.NewClient(
		vc.WithTransport(&httpwares.TestDataTransport{
			Status:      status,
			Filename:    filename,
			ContentType: "application/json",
		}),
		vc.WithAPIKey("fooKey"))
}
