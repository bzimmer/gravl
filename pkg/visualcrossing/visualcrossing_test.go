package visualcrossing_test

import (
	"context"
	"net/http"
	"testing"

	vc "github.com/bzimmer/gravl/pkg/visualcrossing"
	"github.com/bzimmer/transport"

	"github.com/stretchr/testify/assert"
)

func newClient(status int, filename string) (*vc.Client, error) {
	return vc.NewClient(
		vc.WithTransport(&transport.TestDataTransport{
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

	a.Equal(16, len(fcst.Period.Conditions))

	conditions := fcst.Period.Conditions
	cond := conditions[len(conditions)-1]
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

	fault := err.(*vc.Fault)
	a.Equal(106, fault.ErrorCode)
	a.Equal("No session found with id 'null'. The session may have expired", fault.Error())
}
