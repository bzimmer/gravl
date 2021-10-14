//go:build integration

package cyclinganalytics_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bzimmer/gravl/pkg/internal"
)

func TestActivityIntegration(t *testing.T) {
	t.Skip()
	suite.Run(t, &internal.ActivityTestSuite{
		Name:       "ca",
		Encodings:  []string{"gpx", "named"},
		Routes:     false,
		Upload:     true,
		StreamSets: true,
	})
}
