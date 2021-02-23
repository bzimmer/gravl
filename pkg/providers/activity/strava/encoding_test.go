package strava_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGPXFromStreams(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	ctx := context.Background()
	client, err := newClient(http.StatusOK, "streams_four.json")
	a.NoError(err)

	streams, err := client.Activity.Streams(ctx, 154504250376, "latlng", "altitude")
	a.NoError(err)
	a.NotNil(streams)
	a.NotNil(streams.LatLng)
	a.NotNil(streams.Elevation)
	a.Equal(int64(154504250376), streams.ActivityID)
	a.Equal(2712, len(streams.LatLng.Data))

	gpx, err := streams.GPX()
	a.NoError(err)
	a.NotNil(gpx)
	a.Equal(2712, len(gpx.Trk[0].TrkSeg[0].TrkPt))
	a.Equal(0, len(gpx.Rte))
}

func TestGPXFromRoute(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	client, err := newClient(http.StatusOK, "route.json")
	a.NoError(err)
	ctx := context.Background()
	rte, err := client.Route.Route(ctx, 26587226)
	a.NoError(err)
	a.NotNil(rte)
	a.Equal(int64(26587226), rte.ID)

	gpx, err := rte.GPX()
	a.NoError(err)
	a.NotNil(gpx)
	a.Equal(2076, len(gpx.Rte[0].RtePt))
	a.Equal(0, len(gpx.Trk))
}

func TestGPXFromActivity(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	client, err := newClient(http.StatusOK, "activity.json")
	a.NoError(err)
	ctx := context.Background()
	act, err := client.Activity.Activity(ctx, 154504250376823)
	a.NoError(err)
	a.NotNil(act)

	gpx, err := act.GPX()
	a.Error(err)
	a.Nil(gpx)

	client, err = newClient(http.StatusOK, "activity_with_polyline.json")
	a.NoError(err)
	ctx = context.Background()
	act, err = client.Activity.Activity(ctx, 154504250376823)
	a.NoError(err)
	a.NotNil(act)

	gpx, err = act.GPX()
	a.NoError(err)
	a.NotNil(gpx)
	a.Equal(7, len(gpx.Trk[0].TrkSeg[0].TrkPt))
	a.Equal(0, len(gpx.Rte))
}
