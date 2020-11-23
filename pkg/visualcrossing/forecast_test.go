package visualcrossing_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bzimmer/gravl/pkg/visualcrossing"
	"github.com/stretchr/testify/assert"
)

func Test_ForecastSuccess(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, err := newClient(http.StatusOK, "forecast.json")
	a.NoError(err)
	a.NotNil(c)

	ctx := context.Background()
	fcst, err := c.Forecast.Forecast(ctx)
	a.NoError(err)
	a.NotNil(fcst)

	loc := fcst.Locations[0]
	a.Equal(16, len(loc.ForecastConditions))

	conditions := loc.ForecastConditions
	cond := conditions[len(conditions)-1]
	a.Equal(32.1, cond.WindChill)

	f, err := fcst.Forecast()
	a.NoError(err)
	a.NotNil(f)
}

func Test_ForecastError(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, err := newClient(http.StatusOK, "error.json")
	a.NoError(err)
	a.NotNil(c)

	ctx := context.Background()
	fcst, err := c.Forecast.Forecast(ctx)
	a.Error(err)
	a.Nil(fcst)

	fault := err.(*visualcrossing.Fault)
	a.Equal(106, fault.ErrorCode)
	a.Equal("No session found with id 'null'. The session may have expired", fault.Error())
}
