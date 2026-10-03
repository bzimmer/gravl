package strava_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	api "github.com/bzimmer/activity/strava"
	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/activity/strava"
	"github.com/bzimmer/gravl/internal"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestRateLimitTransportError(t *testing.T) {
	a := assert.New(t)
	failure := errors.New("connection refused")
	transport := &strava.RateLimitTransport{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, failure }),
	}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://localhost/athlete", nil)
	a.NoError(err)
	res, err := transport.RoundTrip(req)
	a.ErrorIs(err, failure)
	a.Nil(res)
	a.Nil(transport.RateLimits())
}

func TestParseRateLimitsMalformed(t *testing.T) {
	a := assert.New(t)
	for name, headers := range map[string]map[string]string{
		"missing":          {},
		"no comma":         {"X-RateLimit-Limit": "600", "X-RateLimit-Usage": "1,2"},
		"bad 15 minute":    {"X-RateLimit-Limit": "x,6000", "X-RateLimit-Usage": "1,2"},
		"bad daily":        {"X-RateLimit-Limit": "600,x", "X-RateLimit-Usage": "1,2"},
		"bad usage":        {"X-RateLimit-Limit": "600,6000", "X-RateLimit-Usage": "1"},
		"usage without it": {"X-RateLimit-Limit": "600,6000"},
	} {
		header := http.Header{}
		for k, v := range headers {
			header.Set(k, v)
		}
		limits, ok := strava.ParseRateLimits(header)
		a.False(ok, name)
		a.Nil(limits, name)
	}

	// malformed read headers drop only the read limits
	header := http.Header{}
	header.Set("X-RateLimit-Limit", "600,6000")
	header.Set("X-RateLimit-Usage", "1,2")
	header.Set("X-ReadRateLimit-Limit", "300,3000")
	header.Set("X-ReadRateLimit-Usage", "nope")
	limits, ok := strava.ParseRateLimits(header)
	a.True(ok)
	a.Equal(600, limits.Overall.Limit.FifteenMinute)
	a.Nil(limits.Read)
}

func TestRateLimitsFailures(t *testing.T) {
	a := assert.New(t)
	tests := []struct {
		name    string
		handler http.HandlerFunc
		err     string
	}{
		{
			name: "athlete request fails",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				a.NoError(json.NewEncoder(w).Encode(map[string]any{"message": "boom"}))
			},
			err: "boom",
		},
		{
			name: "no headers",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				a.NoError(json.NewEncoder(w).Encode(&api.Athlete{}))
			},
			err: "no rate limit headers",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// a previous response's headers must not stand in for this one's
			strava.RateLimitRecorder().Reset()
			ok := http.NewServeMux()
			ok.HandleFunc("/athlete", func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("X-RateLimit-Limit", "600,6000")
				w.Header().Set("X-RateLimit-Usage", "1,2")
				a.NoError(json.NewEncoder(w).Encode(&api.Athlete{}))
			})
			internal.Run(t, &internal.Harness{Name: "prime", Args: []string{"gravl", "strava", "ratelimits"}}, ok, command)
			a.NotNil(strava.RateLimitRecorder().RateLimits())

			mux := http.NewServeMux()
			mux.HandleFunc("/athlete", tc.handler)
			tt := &internal.Harness{Name: tc.name, Args: []string{"gravl", "strava", "ratelimits"}, Err: tc.err}
			internal.Run(t, tt, mux, command)
		})
	}
}
