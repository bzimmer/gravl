package visualcrossing_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bzimmer/gravl/pkg/common"
	vc "github.com/bzimmer/gravl/pkg/visualcrossing"

	"github.com/stretchr/testify/assert"
)

func newClient(status int, filename string) (*vc.Client, error) {
	return vc.NewClient(
		vc.WithTransport(&common.TestDataTransport{
			Status:      status,
			Filename:    filename,
			ContentType: "application/json",
		}),
		vc.WithAPIKey("fooKey"))
}

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
	a.Equal(1, fcst.QueryCost)
	a.Equal(1, len(fcst.Locations))

	loc := fcst.Locations[0]
	a.Equal(16, len(loc.Conditions))

	cond := loc.Conditions[len(loc.Conditions)-1]
	a.Equal(32.1, cond.WindChill)
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

	fault := err.(vc.Fault)
	a.Equal(106, fault.ErrorCode)
	a.Equal("No session found with id 'null'. The session may have expired", fault.Error())
}
