package rwgps_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/common/route"
)

func contextNil() context.Context {
	return nil
}

func Test_Trip(t *testing.T) { // nolint
	t.Parallel()
	a := assert.New(t)

	c, err := newClient(http.StatusOK, "rwgps_trip_94.json")
	a.NoError(err)
	a.NotNil(c)

	ctx := context.Background()
	rte, err := c.Trips.Trip(ctx, 94)
	a.NoError(err)
	a.NotNil(rte)
	a.Equal("94", rte.ID)
	a.Equal(route.Activity, rte.Origin)
	a.Equal(1465, len(rte.Coordinates))

	rte, err = c.Trips.Trip(contextNil(), 94)
	a.Error(err)
	a.Nil(rte)
}
func Test_Route(t *testing.T) { // nolint
	t.Parallel()
	a := assert.New(t)

	c, err := newClient(http.StatusOK, "rwgps_route_141014.json")
	a.NoError(err)
	a.NotNil(c)

	ctx := context.Background()
	rte, err := c.Trips.Route(ctx, 141014)
	a.NoError(err)
	a.NotNil(rte)
	a.Equal("141014", rte.ID)
	a.Equal(route.Planned, rte.Origin)
	a.Equal(1154, len(rte.Coordinates))

	// keep the linter quiet by using a function to return nil
	rte, err = c.Trips.Route(contextNil(), 141014)
	a.Error(err)
	a.Nil(rte)
}
