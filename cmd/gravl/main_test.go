package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	api "github.com/bzimmer/activity/strava"
)

type fakeFaultError struct {
	code int
}

func (f *fakeFaultError) Error() string       { return "fake fault" }
func (f *fakeFaultError) HTTPStatusCode() int { return f.code }

func TestExitCode(t *testing.T) {
	for _, tt := range []struct {
		name string
		err  error
		code int
	}{
		{name: "nil error", err: nil, code: 1},
		{name: "plain error", err: errors.New("boom"), code: 1},
		{name: "non-rate-limit fault", err: &fakeFaultError{code: http.StatusInternalServerError}, code: 1},
		{name: "rate limited fault", err: &fakeFaultError{code: http.StatusTooManyRequests}, code: exitTempFail},
		{
			name: "wrapped rate limited fault",
			err:  fmt.Errorf("wrap: %w", &fakeFaultError{code: http.StatusTooManyRequests}),
			code: exitTempFail,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.code, exitCode(tt.err))
		})
	}
}

// TestExitCodeStravaRateLimit proves the full chain against a real *strava.Fault,
// not the hand-rolled fakeFaultError above: a genuine 429 HTTP response, decoded
// by github.com/bzimmer/activity's client into strava.Fault, fed into exitCode().
//
// This stops short of a real subprocess run of the gravl binary because
// strava.Before (activity/strava/strava.go) has no base-URL override -- it always
// refreshes against the live Strava OAuth endpoint -- so "gravl strava activity"
// can't be pointed at a fake server without a production code change.
func TestExitCodeStravaRateLimit(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/athlete", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"message": "Rate Limit Exceeded"}`))
	})
	svr := httptest.NewServer(mux)
	defer svr.Close()

	client, err := api.NewClient(
		api.WithBaseURL(svr.URL),
		api.WithTokenCredentials("token", "token", time.Now().Add(time.Hour)))
	assert.NoError(t, err)

	_, err = client.Athlete.Athlete(context.Background())
	assert.Error(t, err)

	assert.Equal(t, exitTempFail, exitCode(err))
}
