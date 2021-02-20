package strava_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/bzimmer/httpwares"
	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/providers/activity"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

func readall(ctx context.Context, acts <-chan *strava.ActivityResult) ([]*strava.Activity, error) {
	var activities []*strava.Activity
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case res, ok := <-acts:
			if !ok {
				// the channel is closed, return the activities
				return activities, nil
			}
			if res.Err != nil {
				return nil, res.Err
			}
			activities = append(activities, res.Activity)
		}
	}
}

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

	acts, err := readall(ctx, client.Activity.Activities(ctx, activity.Pagination{}))
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
	acts, err := readall(ctx, client.Activity.Activities(ctx, activity.Pagination{Total: 325}))
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

	t.Run("total, start, and count", func(t *testing.T) {
		// success: the requested number of activities because count/pagesize == 1
		acts, err := readall(ctx, client.Activity.Activities(ctx, activity.Pagination{Total: 127, Start: 0, Count: 1}))
		a.NoError(err)
		a.NotNil(acts)
		a.Equal(127, len(acts))
	})

	t.Run("total and start", func(t *testing.T) {
		// success: the requested number of activities is exceeded because count/pagesize not specified
		x := 234
		acts, err := readall(ctx, client.Activity.Activities(ctx, activity.Pagination{Total: x, Start: 0}))
		a.NoError(err)
		a.NotNil(acts)
		a.Equal(x, len(acts))
	})

	t.Run("total and start less than PageSize", func(t *testing.T) {
		// success: the requested number of activities because count/pagesize <= strava.PageSize
		a.True(27 < strava.PageSize)
		acts, err := readall(ctx, client.Activity.Activities(ctx, activity.Pagination{Total: 27, Start: 0}))
		a.NoError(err)
		a.NotNil(acts)
		a.Equal(27, len(acts))
	})

	t.Run("different Count values", func(t *testing.T) {
		count := strava.PageSize + 100
		for _, x := range []int{27, 350, strava.PageSize} {
			acts, err := readall(ctx, client.Activity.Activities(ctx, activity.Pagination{Total: x, Start: 0, Count: count}))
			a.NoError(err)
			a.NotNil(acts)
			a.Equal(x, len(acts))
		}
	})

	t.Run("negative total", func(t *testing.T) {
		acts, err := readall(ctx, client.Activity.Activities(ctx, activity.Pagination{Total: -1}))
		a.Error(err)
		a.Nil(acts)
	})
}

func TestStreams(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("four", func(t *testing.T) {
		ctx := context.Background()
		client, err := newClient(http.StatusOK, "streams_four.json")
		a.NoError(err)
		sms, err := client.Activity.Streams(ctx, 154504250376, "latlng", "altitude", "distance")
		a.NoError(err)
		a.NotNil(sms)
		a.NotNil(sms.LatLng)
		a.NotNil(sms.Elevation)
		a.NotNil(sms.Distance)
	})

	t.Run("two", func(t *testing.T) {
		ctx := context.Background()
		client, err := newClient(http.StatusOK, "streams_two.json")
		a.NoError(err)
		sms, err := client.Activity.Streams(ctx, 154504250376, "latlng", "altitude")
		a.NoError(err)
		a.NotNil(sms)
		a.NotNil(sms.LatLng)
		a.NotNil(sms.Elevation)
	})
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

	t.Run("timeout lt sleep => failure", func(t *testing.T) {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Millisecond*15)
		defer cancel()
		act, err := client.Activity.Activity(ctx, 154504250376823)
		a.Error(err)
		a.Nil(act)
	})

	t.Run("timeout gt sleep => success", func(t *testing.T) {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Millisecond*120)
		defer cancel()
		act, err := client.Activity.Activity(ctx, 154504250376823)
		a.NoError(err)
		a.NotNil(act)
		a.Equal(int64(154504250376823), act.ID)
	})
}
