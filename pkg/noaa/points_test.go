package noaa_test

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
	var requester = func(req *http.Request) error {
		a.Equal("https://api.weather.gov/points/48.03,-121.44/forecast", req.URL.String())
		return nil
	}
	c, err := newClienter(http.StatusOK, "barlow_pass_forecast_daily.json", requester, nil)
	a.NoError(err)
	a.NotNil(c)
	f, err := c.Points.Forecast(b, -121.4440005, 48.0264959)
	a.NoError(err)
	a.NotNil(f)
	a.Equal(14, len(f.Properties.Periods))

	// hourly
	requester = func(req *http.Request) error {
		a.Equal("https://api.weather.gov/points/48.03,-121.44/forecast/hourly", req.URL.String())
		return nil
	}
	c, err = newClienter(http.StatusOK, "barlow_pass_forecast_hourly.json", requester, nil)
	a.NoError(err)
	a.NotNil(c)
	f, err = c.Points.ForecastHourly(b, -121.4440005, 48.0264959)
	a.NoError(err)
	a.NotNil(f)
	a.Equal(156, len(f.Properties.Periods))

	// failure
	requester = func(req *http.Request) error {
		a.Equal("https://api.weather.gov/points/2.03,-121.44/forecast", req.URL.String())
		return nil
	}
	c, err = newClienter(http.StatusNotFound, "unavailable_forecast.json", requester, nil)
	a.NoError(err)
	a.NotNil(c)
	f, err = c.Points.Forecast(b, -121.444, 2.0265)
	a.Error(err)
	a.Nil(f)
	a.Equal("Unable to provide data for requested point 2.0265,-121.444", err.Error())
}

func Test_Gridpoint(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var requester = func(req *http.Request) error {
		a.Equal("https://api.weather.gov/points/48.03,-121.44", req.URL.String())
		return nil
	}
	c, err := newClienter(http.StatusOK, "barlow_pass_gridpoint.json", requester, nil)
	a.NoError(err)
	a.NotNil(c)

	b := context.Background()
	p, err := c.Points.GridPoint(b, -121.4440005, 48.0264959)
	a.NoError(err)
	a.NotNil(p)
	a.Equal("SEW", p.Properties.GridID)
	a.Equal(156, p.Properties.GridX)
	a.Equal(81, p.Properties.GridY)
}
