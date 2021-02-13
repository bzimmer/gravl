package openweather_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/commands/internal"
)

func TestForecastIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	a := assert.New(t)

	c := internal.Gravl("ow", "forecast", "47.62", "-121.52")
	<-c.Start()
	a.True(c.Success())
}
