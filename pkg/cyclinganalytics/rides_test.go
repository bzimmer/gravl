package cyclinganalytics_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bzimmer/gravl/pkg/cyclinganalytics"
	"github.com/stretchr/testify/assert"
)

func TestRides(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	client, err := newClient(http.StatusOK, "me-rides.json")
	a.NoError(err)
	ctx := context.Background()
	rides, err := client.Rides.Rides(ctx, cyclinganalytics.Me)
	a.NoError(err)
	a.NotNil(rides)
	a.Equal(2, len(rides))
}

func TestRide(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	client, err := newClient(http.StatusOK, "ride.json")
	a.NoError(err)
	ctx := context.Background()
	opts := cyclinganalytics.RideOptions{
		Streams: []string{"latitude", "longitude", "elevation"},
	}
	ride, err := client.Rides.Ride(ctx, 175334338355, opts)
	a.NoError(err)
	a.NotNil(ride)
	a.NotNil(ride.Streams)
	a.Equal(27154, len(ride.Streams.Elevation))
	gears := ride.Streams.Gears
	a.NotNil(gears)
	a.Equal(813, len(gears.Shifts))
}

func TestRideForbidden(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	client, err := newClient(http.StatusForbidden, "error.json")
	a.NoError(err)
	ctx := context.Background()
	ride, err := client.Rides.Ride(ctx, 175334338355, cyclinganalytics.RideOptions{})
	a.Error(err)
	a.Nil(ride)
}
