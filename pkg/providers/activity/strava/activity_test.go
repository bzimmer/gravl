package strava_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/bzimmer/httpwares"
	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/providers/activity"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

func TestActivity(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	client, err := newClient(http.StatusOK, "activity.json")
	a.NoError(err)
	ctx := context.Background()
	act, err := client.Activity.Activity(ctx, 154504250376823)
	a.NoError(err)
	a.NotNil(act)
	a.Equal(int64(154504250376823), act.ID)
}

func TestActivities(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	ctx := context.Background()
	client, err := strava.NewClient(
		strava.WithTransport(&ManyTransport{
			Filename: "testdata/activity.json",
			Total:    2,
		}),
		strava.WithTokenCredentials("fooKey", "barToken", time.Time{}))
	a.NoError(err)

	acts, err := client.Activity.Activities(ctx, activity.Pagination{})
	a.NoError(err)
	a.Equal(2, len(acts))
}

type F struct {
	n int
}

func (f *F) X(res *http.Response) error {
	if f.n == 1 {
		// on the second iteration return an empty body signaling no more activities exist
		res.ContentLength = int64(0)
		res.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))
	}
	f.n++
	return nil
}

func TestActivitiesRequestedGTAvailable(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	client, err := newClienter(http.StatusOK, "activities.json", nil, (&F{}).X)
	a.NoError(err)
	ctx := context.Background()
	acts, err := client.Activity.Activities(ctx, activity.Pagination{Total: 325})
	a.NoError(err)
	a.Equal(2, len(acts))
}

func TestActivitiesMany(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	ctx := context.Background()
	client, err := strava.NewClient(
		strava.WithTransport(&ManyTransport{
			Filename: "testdata/activity.json",
		}),
		strava.WithTokenCredentials("fooKey", "barToken", time.Time{}))
	a.NoError(err)

	// test total, start, and count
	// success: the requested number of activities because count/pagesize == 1
	acts, err := client.Activity.Activities(ctx, activity.Pagination{Total: 127, Start: 0, Count: 1})
	a.NoError(err)
	a.NotNil(acts)
	a.Equal(127, len(acts))

	// test total and start
	// success: the requested number of activities is exceeded because count/pagesize not specified
	x := 234
	n := int(math.Floor(float64(x)/strava.PageSize)*strava.PageSize + strava.PageSize)
	acts, err = client.Activity.Activities(ctx, activity.Pagination{Total: x, Start: 0})
	a.NoError(err)
	a.NotNil(acts)
	a.Equal(n, len(acts))

	// test total and start less than PageSize
	// success: the requested number of activities because count/pagesize <= strava.PageSize
	a.True(27 < strava.PageSize)
	acts, err = client.Activity.Activities(ctx, activity.Pagination{Total: 27, Start: 0})
	a.NoError(err)
	a.NotNil(acts)
	a.Equal(27, len(acts))

	// test different Count values
	count := strava.PageSize + 100
	for _, x = range []int{27, 350, strava.PageSize} {
		acts, err = client.Activity.Activities(ctx, activity.Pagination{Total: x, Start: 0, Count: count})
		a.NoError(err)
		a.NotNil(acts)

		n = x
		if x > strava.PageSize {
			n = int(math.Floor(float64(x)/strava.PageSize)*strava.PageSize + strava.PageSize)
		}
		a.Equal(n, len(acts))
	}

	// negative test
	acts, err = client.Activity.Activities(ctx, activity.Pagination{Total: -1})
	a.Error(err)
	a.Nil(acts)
}

func TestStreams(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	ctx := context.Background()
	client, err := newClient(http.StatusOK, "streams_four.json")
	a.NoError(err)

	sms, err := client.Activity.Streams(ctx, 154504250376, "latlng", "altitude", "distance")
	a.NoError(err)
	a.NotNil(sms)
	a.NotNil(sms.LatLng)
	a.NotNil(sms.Elevation)
	a.NotNil(sms.Distance)

	client, err = newClient(http.StatusOK, "streams_two.json")
	a.NoError(err)
	sms, err = client.Activity.Streams(ctx, 154504250376, "latlng", "altitude")
	a.NoError(err)
	a.NotNil(sms)
	a.NotNil(sms.LatLng)
	a.NotNil(sms.Elevation)
}

func TestRouteFromStreams(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	ctx := context.Background()
	client, err := newClient(http.StatusOK, "streams_four.json")
	a.NoError(err)

	sms, err := client.Activity.Streams(ctx, 154504250376, "latlng", "altitude")
	a.NoError(err)
	a.NotNil(sms)
	a.NotNil(sms.LatLng)
	a.NotNil(sms.Elevation)
	a.Equal(int64(154504250376), sms.ActivityID)
	a.Equal(2712, len(sms.LatLng.Data))
}

func TestTimeout(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	client, err := strava.NewClient(
		strava.WithTokenCredentials("fooKey", "barToken", time.Time{}),
		strava.WithTransport(&httpwares.SleepingTransport{
			Duration: time.Millisecond * 30,
			Transport: &httpwares.TestDataTransport{
				Status:      http.StatusOK,
				Filename:    "activity.json",
				ContentType: "application/json",
			}}))
	a.NoError(err)
	a.NotNil(client)

	// timeout lt sleep => failure
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*15)
	defer cancel()

	act, err := client.Activity.Activity(ctx, 154504250376823)
	a.Error(err)
	a.Nil(act)

	// timeout gt sleep => success
	ctx = context.Background()
	ctx, cancel = context.WithTimeout(ctx, time.Millisecond*45)
	defer cancel()

	act, err = client.Activity.Activity(ctx, 154504250376823)
	a.NoError(err)
	a.NotNil(act)
	a.Equal(int64(154504250376823), act.ID)
}
