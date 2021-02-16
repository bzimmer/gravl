package noaa_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bzimmer/gravl/pkg/commands/wx/internal"
)

func TestForecastIntegration(t *testing.T) {
	t.Skipf("skipping noaa's seriously unreliable api tests")
	suite.Run(t, &internal.ForecastTestSuite{
		Name: "noaa",
	})
}
