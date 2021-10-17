//go:build integration

package strava_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bzimmer/gravl/pkg/internal"
)

func TestActivityIntegration(t *testing.T) {
	suite.Run(t, &internal.ActivityTestSuite{
		Name:       "strava",
		Encodings:  []string{"gpx", "named"},
		Routes:     true,
		Upload:     true,
		StreamSets: true,
	})
}
