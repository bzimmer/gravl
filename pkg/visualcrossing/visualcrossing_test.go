package visualcrossing_test

import (
	"github.com/bzimmer/gravl/pkg/common"
	vc "github.com/bzimmer/gravl/pkg/visualcrossing"
)

func newClient(status int, filename string) (*vc.Client, error) {
	return vc.NewClient(
		vc.WithTransport(&common.TestDataTransport{
			Status:      status,
			Filename:    filename,
			ContentType: "application/json",
		}),
		vc.WithAPIKey("fooKey"))
}
