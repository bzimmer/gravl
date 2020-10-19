package noaa

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Points_Forecast(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	b := context.Background()

	// daily
	c, err := newClient(http.StatusOK, "barlow_pass_forecast_daily.json")
	a.NoError(err)
	a.NotNil(c)
	f, err := c.Points.Forecast(b, 48.0264959, -121.4440005)
	a.NoError(err)
	a.NotNil(f)
	a.Equal(14, len(f.Properties.Periods))

	// hourly
	c, err = newClient(http.StatusOK, "barlow_pass_forecast_hourly.json")
	a.NoError(err)
	a.NotNil(c)
	f, err = c.Points.ForecastHourly(b, 48.0264959, -121.4440005)
	a.NoError(err)
	a.NotNil(f)
	a.Equal(156, len(f.Properties.Periods))
	a.Equal("2020-10-23T16:00:00-07:00", f.Properties.Periods[118].StartTime)

	// failure
	c, err = newClient(http.StatusNotFound, "unavailable_forecast.json")
	a.NoError(err)
	a.NotNil(c)
	f, err = c.Points.Forecast(b, 2.0265, -121.444)
	a.Error(err)
	a.Nil(f)
	a.Equal("Unable to provide data for requested point 2.0265,-121.444", err.Error())
}

func Test_Gridpoint(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, err := newClient(http.StatusOK, "barlow_pass_gridpoint.json")
	a.NoError(err)
	a.NotNil(c)

	b := context.Background()
	p, err := c.Points.GridPoint(b, 48.0264959, -121.4440005)
	a.NoError(err)
	a.NotNil(p)
	a.Equal("SEW", p.Properties.GridID)
	a.Equal(156, p.Properties.GridX)
	a.Equal(81, p.Properties.GridY)
}