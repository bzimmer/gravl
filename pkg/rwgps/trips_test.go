package rwgps_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/wta/pkg/common"
	rw "github.com/bzimmer/wta/pkg/rwgps"
)

func newClient(status int, filename string) (*rw.Client, error) {
	return rw.NewClient(
		rw.WithTransport(&common.TestDataTransport{
			Status:      status,
			Filename:    filename,
			ContentType: "application/json",
		}),
	)
}

func Test_Trip(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, err := newClient(http.StatusOK, "rwgps_trip_94.json")
	a.NoError(err)
	a.NotNil(c)

	ctx := context.Background()
	fc, err := c.Trips.Trip(ctx, 94)
	a.NoError(err)
	a.NotNil(fc)
	a.Equal(94, fc.Features[0].ID)
	a.Equal("trip", fc.Features[0].Properties["type"])
	a.True(fc.Features[0].Geometry.IsLineString())
	a.Equal(1465, len(fc.Features[0].Geometry.LineString))

	fc, err = c.Trips.Trip(nil, 94)
	a.Error(err)
	a.Nil(fc)
}

func Test_Route(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, err := newClient(http.StatusOK, "rwgps_route_141014.json")
	a.NoError(err)
	a.NotNil(c)

	ctx := context.Background()
	fc, err := c.Trips.Route(ctx, 141014)
	a.NoError(err)
	a.NotNil(fc)
	a.Equal(141014, fc.Features[0].ID)
	a.Equal("route", fc.Features[0].Properties["type"])
	a.Equal(1, len(fc.Features))
	a.True(fc.Features[0].Geometry.IsLineString())
	a.Equal(1154, len(fc.Features[0].Geometry.LineString))

	fc, err = c.Trips.Route(nil, 141014)
	a.Error(err)
	a.Nil(fc)
}
