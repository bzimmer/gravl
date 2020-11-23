package openweather_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/bzimmer/gravl/pkg/openweather"

	"github.com/stretchr/testify/assert"
)

func Test_ForecastSuccess(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, err := newClient(http.StatusOK, "onecall.json")
	a.NoError(err)
	a.NotNil(c)

	ctx := context.Background()
	fcst, err := c.Forecast.Forecast(ctx)
	a.NoError(err)
	a.NotNil(fcst)

	hourly := fcst.Hourly
	a.Equal(48, len(hourly))

	f, err := fcst.Forecast()
	a.Error(err)
	a.Nil(f)
}

func TestWithUnits(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	v := &url.Values{}
	tests := map[string]openweather.Units{
		"metric":   openweather.UnitsMetric,
		"imperial": openweather.UnitsImperial,
		"standard": openweather.UnitsStandard,
	}

	for key, value := range tests {
		err := openweather.WithUnits(value)(v)
		a.NoError(err)
		a.Equal(key, v.Get("units"))
	}
}

func TestWithLocation(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	v := &url.Values{}
	err := openweather.WithCoordinates(-122.05, 48.10)(v)
	a.NoError(err)
	a.Equal("-122.0500", v.Get("lon"))
	a.Equal("48.1000", v.Get("lat"))
}
