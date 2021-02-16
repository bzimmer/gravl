package cyclinganalytics_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bzimmer/gravl/pkg/commands/activity/internal"
)

func TestActivityIntegration(t *testing.T) {
	suite.Run(t, &internal.ActivityTestSuite{
		Name:       "ca",
		Encodings:  []string{"gpx"},
		SkipRoutes: true,
	})
}
