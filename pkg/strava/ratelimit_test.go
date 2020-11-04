package strava_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/strava"
)

func Test_RateLimit(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	r := &strava.RateLimit{
		LimitWindow: 100,
		LimitDaily:  500,
		UsageWindow: 50,
		UsageDaily:  250,
	}
	a.Equal(50, r.PercentWindow())
	a.Equal(50, r.PercentDaily())
}

func Test_RateLimitExceeded(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	client, err := newClienter(
		http.StatusTooManyRequests,
		"exceeded_rate_limit.json",
		nil,
		func(res *http.Response) error {
			res.Header.Add(strava.HeaderRateLimit, "600,30000")
			res.Header.Add(strava.HeaderRateUsage, "601,30100")
			return nil
		})
	strava.WithRateLimitThrottling()(client)
	a.NoError(err)

	// call the first time seed the client with the rate limit response
	//  this will result in a too many requests error
	ctx := context.Background()
	sts, err := client.Athlete.Stats(ctx, 88273)
	a.Nil(sts)
	a.Error(err)
	a.Error(err.(*strava.Fault))
	a.Equal("exceeded rate limit", err.(*strava.Fault).Message)

	// the second call will fail not with the Fault but a RateLimitError
	//  (wrapped by url.Error) which can be inspected and used to throttle
	sts, err = client.Athlete.Stats(ctx, 88273)
	a.Nil(sts)
	a.Error(err)
	er := err.(*url.Error).Unwrap()
	a.Error(er.(*strava.RateLimitError))
	r := (er.(*strava.RateLimitError)).RateLimit
	a.Equal(30000, r.LimitDaily)
	a.Equal(601, r.UsageWindow)
}

func Test_RateLimitTransport(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	client, err := newClienter(
		http.StatusOK,
		"athlete_stats.json",
		nil,
		func(res *http.Response) error {
			res.Header.Add(strava.HeaderRateLimit, "600,30000")
			res.Header.Add(strava.HeaderRateUsage, "314,27536")
			return nil
		})
	a.NoError(err)
	strava.WithRateLimitThrottling()(client)

	ctx := context.Background()
	sts, err := client.Athlete.Stats(ctx, 88273)
	a.NotNil(sts)
	a.Nil(err)
}
