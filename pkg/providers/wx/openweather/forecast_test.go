package openweather_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twpayne/go-geom"

	"github.com/bzimmer/gravl/pkg/providers/wx"
)

func Test_ForecastSuccess(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, err := newClient(http.StatusOK, "onecall.json")
	a.NoError(err)
	a.NotNil(c)

	ctx := context.Background()
	opts := wx.ForecastOptions{
		Point: geom.NewPointFlat(geom.XY, []float64{-118.8, 48.2})}
	fcst, err := c.Forecast.Forecast(ctx, opts)
	a.NoError(err)
	a.NotNil(fcst)

	hourly := fcst.Hourly
	a.Equal(48, len(hourly))

	f, err := fcst.Forecast()
	a.Error(err)
	a.Nil(f)
}
