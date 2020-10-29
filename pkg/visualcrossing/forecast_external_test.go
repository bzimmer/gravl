package visualcrossing_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Forecast(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, err := newClient(http.StatusOK, "98110_forecast_array.json")
	a.NoError(err)
	a.NotNil(c)

	ctx := context.Background()
	fcst, err := c.Forecast.Forecast(ctx)
	a.NoError(err)
	a.NotNil(fcst)
	a.Equal(1, fcst.QueryCost)
	a.Equal(1, len(fcst.Locations))

	loc := fcst.Locations[0]
	a.Equal(16, len(loc.Conditions))

	cond := loc.Conditions[len(loc.Conditions)-1]
	a.Equal(32.1, cond.WindChill)
}
