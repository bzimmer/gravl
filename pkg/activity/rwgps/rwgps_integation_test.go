//go:build integration

package rwgps_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bzimmer/gravl/pkg/internal"
)

func TestActivityIntegration(t *testing.T) {
	suite.Run(t, &internal.ActivityTestSuite{
		Name:      "rwgps",
		Encodings: []string{"gpx", "geojson", "named"},
		Routes:    true,
		Upload:    true,
	})
}
