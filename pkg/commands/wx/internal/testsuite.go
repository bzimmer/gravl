package internal

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bzimmer/gravl/pkg/commands/internal"
)

// ForecastTestSuite for testing forecasting services
type ForecastTestSuite struct {
	suite.Suite
	Name string
}

func (s *ForecastTestSuite) SetupSuite() {
	if testing.Short() {
		s.T().Skip("skipping integration test suite")
		return
	}
}

func (s *ForecastTestSuite) BeforeTest(suiteName, testName string) {
	s.T().Parallel()
}

func (s *ForecastTestSuite) TestForecast() {
	a := s.Assert()
	c := internal.Gravl(s.Name, "forecast", "47.62", "-121.52")
	<-c.Start()
	a.True(c.Success())
}
