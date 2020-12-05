package rwgps_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bzimmer/gravl/pkg/rwgps"
	"github.com/stretchr/testify/assert"
)

func contextNil() context.Context {
	return nil
}

func Test_Trip(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, err := newClient(http.StatusOK, "rwgps_trip_94.json")
	a.NoError(err)
	a.NotNil(c)

	ctx := context.Background()
	trk, err := c.Trips.Trip(ctx, 94)
	a.NoError(err)
	a.NotNil(trk)
	a.Equal(int64(94), trk.ID)
	a.Equal(rwgps.OriginTrip, trk.Origin)
	a.Equal(1465, len(trk.TrackPoints))
	a.Equal("OriginTrip", trk.Origin.String())

	trk, err = c.Trips.Trip(contextNil(), 94)
	a.Error(err)
	a.Nil(trk)
}

func Test_Route(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, err := newClient(http.StatusOK, "rwgps_route_141014.json")
	a.NoError(err)
	a.NotNil(c)

	ctx := context.Background()
	rte, err := c.Trips.Route(ctx, 141014)
	a.NoError(err)
	a.NotNil(rte)
	a.Equal(1154, len(rte.TrackPoints))
	a.Equal(int64(141014), rte.ID)
	a.Equal(rwgps.OriginRoute, rte.Origin)

	// trk, err := rte.Track()
	// a.NoError(err)
	// a.NotNil(trk)
	// a.Equal(1154, len(trk.Coordinates))

	// keep the linter quiet by using a function to return nil
	rte, err = c.Trips.Route(contextNil(), 141014)
	a.Error(err)
	a.Nil(rte)
}
