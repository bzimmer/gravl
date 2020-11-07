package strava_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/stretchr/testify/assert"
)

func Test_Activity(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	client, err := newClient(http.StatusOK, "activity.json")
	ctx := context.Background()
	act, err := client.Activity.Activity(ctx, 154504250376823)
	a.NoError(err)
	a.NotNil(act)
	a.Equal(int64(154504250376823), act.ID)
}

func Test_Activities(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	client, err := newClient(http.StatusOK, "activities.json")
	ctx := context.Background()
	acts, err := client.Activity.Activities(ctx)
	a.NoError(err)
	a.Equal(2, len(*acts))
}

func Test_ActivitiesMax(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	client, err := newClient(http.StatusOK, "activities.json")
	ctx := context.Background()
	acts, err := client.Activity.Activities(ctx, 5000)
	a.NoError(err)
	a.Equal(2, len(*acts))
}

type ManyTransport struct{}

func (t *ManyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	q := req.URL.Query()
	n, _ := strconv.Atoi(q.Get("per_page"))

	data, err := ioutil.ReadFile("testdata/activity.json")
	if err != nil {
		return nil, err
	}

	acts := make([]string, 0)
	for i := 0; i < n; i++ {
		acts = append(acts, string(data))
	}

	res := strings.Join(acts, ",")
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString("[" + res + "]")),
		Header:     make(http.Header),
	}, nil
}

func Test_ActivitiesMany(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	ctx := context.Background()
	client, err := strava.NewClient(
		strava.WithTransport(&ManyTransport{}),
		strava.WithAPICredentials("fooKey", "barToken"))
	a.NoError(err)

	// test total, start, and count
	acts, err := client.Activity.Activities(ctx, 352, 0, 1)
	a.NoError(err)
	a.NotNil(acts)
	a.Equal(352, len(*acts))

	// test total and start
	acts, err = client.Activity.Activities(ctx, 234, 0)
	a.NoError(err)
	a.NotNil(acts)
	a.Equal(234, len(*acts))

	// negative test
	acts, err = client.Activity.Activities(ctx, -1)
	a.Error(err)
	a.Nil(acts)

	// test too many varargs
	acts, err = client.Activity.Activities(ctx, 1, 2, 3, 4, 5, 6)
	a.Error(err)
	a.Nil(acts)
}

func Test_Streams(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	ctx := context.Background()
	client, err := newClient(http.StatusOK, "streams.json")
	a.NoError(err)
	fc, err := client.Activity.Streams(ctx, 154504250376, "latlng", "altitude")
	a.NoError(err)
	a.NotNil(fc)
	a.Equal(1, len(fc.Features))

	feature := fc.Features[0]
	a.Equal(int64(154504250376), feature.ID)
	a.True(feature.Geometry.IsLineString())
	a.Equal(2712, len(feature.Geometry.LineString))

	properties := feature.Properties
	streams, ok := (properties["streams"]).(map[string]interface{})
	a.True(ok)
	grades, ok := streams["grade_smooth"]
	a.True(ok)
	a.Equal(2712, len(grades.([]interface{})))
}
