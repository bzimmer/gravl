package noaa_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GridPoints_Forecast(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, err := newClient(http.StatusOK, "barlow_pass_forecast_daily.json")
	a.NoError(err)
	a.NotNil(c)

	b := context.Background()
	f, err := c.GridPoints.Forecast(b, "SEW", 156, 81)
	a.NoError(err)
	a.NotNil(f)
	a.Equal(14, len(f.Properties.Periods))

	c, err = newClient(http.StatusNotFound, "unavailable_forecast.json")
	a.NoError(err)
	a.NotNil(c)
	f, err = c.GridPoints.Forecast(b, "SEW", 156, 81)
	a.Error(err)
	a.Nil(f)
	a.Equal("Unable to provide data for requested point 2.0265,-121.444", err.Error())
}
