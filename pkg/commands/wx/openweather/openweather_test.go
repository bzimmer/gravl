package openweather_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bzimmer/gravl/pkg/commands/wx/internal"
)

func TestForecastIntegration(t *testing.T) {
	suite.Run(t, &internal.ForecastTestSuite{
		Name: "openweather",
	})
}
