package strava_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	api "github.com/bzimmer/activity/strava"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/activity/strava"
	"github.com/bzimmer/gravl/internal"
)

func TestAPIFailures(t *testing.T) {
	tests := []*internal.Harness{
		{Name: "athlete", Args: []string{"gravl", "strava", "athlete"}},
		{Name: "routes", Args: []string{"gravl", "strava", "routes"}},
		{Name: "activities", Args: []string{"gravl", "strava", "activities"}},
		{Name: "one activity", Args: []string{"gravl", "strava", "activity", "12345"}},
		{Name: "several activities", Args: []string{"gravl", "strava", "activity", "1", "2", "3", "4"}},
		{
			Name: "several activities, one at a time",
			Args: []string{"gravl", "strava", "--concurrency", "1", "activity", "1", "2", "3", "4"},
		},
		{Name: "webhook list", Args: []string{"gravl", "strava", "webhook", "list"}},
		{Name: "webhook unsubscribe all", Args: []string{"gravl", "strava", "webhook", "unsubscribe"}},
		{
			Name: "webhook subscribe",
			Args: []string{"gravl", "strava", "webhook", "subscribe", "--url", "https://example.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.Err = "provider is down"
			internal.Run(t, tt, internal.FailingAPI(t), command)
		})
	}
}

func TestRoutesFailures(t *testing.T) {
	a := assert.New(t)
	mux := http.NewServeMux()
	mux.HandleFunc("/athlete", func(w http.ResponseWriter, _ *http.Request) {
		a.NoError(json.NewEncoder(w).Encode(&api.Athlete{ID: 8542982}))
	})
	mux.HandleFunc("/athletes/8542982/routes", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "1" {
			a.NoError(json.NewEncoder(w).Encode([]*api.Route{}))
			return
		}
		a.NoError(json.NewEncoder(w).Encode([]*api.Route{{}}))
	})
	tests := []*internal.Harness{
		{
			Name: "routes request fails",
			Args: []string{"gravl", "strava", "routes"},
			Err:  "provider is down",
		},
		{
			Name:   "encoding fails",
			Args:   []string{"gravl", "strava", "routes"},
			Before: internal.FailEncoding,
			Err:    internal.ErrEncode.Error(),
		},
	}
	for i, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			handler := http.Handler(mux)
			if i == 0 {
				handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/athlete" {
						mux.ServeHTTP(w, r)
						return
					}
					internal.FailingAPI(t).ServeHTTP(w, r)
				})
			}
			internal.Run(t, tt, handler, command)
		})
	}
}

func TestEncodingFailures(t *testing.T) {
	a := assert.New(t)
	mux := http.NewServeMux()
	mux.HandleFunc("/athlete", func(w http.ResponseWriter, _ *http.Request) {
		a.NoError(json.NewEncoder(w).Encode(&api.Athlete{}))
	})
	mux.HandleFunc("/athlete/activities", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "1" {
			a.NoError(json.NewEncoder(w).Encode([]*api.Activity{}))
			return
		}
		a.NoError(json.NewEncoder(w).Encode([]*api.Activity{{Type: "Ride"}}))
	})
	mux.HandleFunc("/push_subscriptions", func(w http.ResponseWriter, _ *http.Request) {
		a.NoError(json.NewEncoder(w).Encode([]*api.WebhookSubscription{{ID: 10}}))
	})
	activities := activityHandler(a)
	mux.Handle("/activities/", activities)

	tests := []*internal.Harness{
		{Name: "athlete", Args: []string{"gravl", "strava", "athlete"}},
		{Name: "activities", Args: []string{"gravl", "strava", "activities"}},
		{Name: "one activity", Args: []string{"gravl", "strava", "activity", "12345"}},
		{Name: "several activities", Args: []string{"gravl", "strava", "activity", "12345", "54321"}},
		{Name: "webhook list", Args: []string{"gravl", "strava", "webhook", "list"}},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.Before = internal.FailEncoding
			tt.Err = internal.ErrEncode.Error()
			internal.Run(t, tt, mux, command)
		})
	}
}

func TestExpressionFailures(t *testing.T) {
	a := assert.New(t)
	mux := http.NewServeMux()
	mux.HandleFunc("/athlete/activities", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "1" {
			a.NoError(json.NewEncoder(w).Encode([]*api.Activity{}))
			return
		}
		a.NoError(json.NewEncoder(w).Encode([]*api.Activity{{Type: "Ride"}}))
	})
	tests := []*internal.Harness{
		{
			Name: "invalid filter",
			Args: []string{"gravl", "strava", "activities", "--filter", ".Type =="},
			Err:  "unexpected token",
		},
		{
			Name: "invalid attribute",
			Args: []string{"gravl", "strava", "activities", "--attribute", ".Type =="},
			Err:  "unexpected token",
		},
		{
			Name: "filter not a boolean",
			Args: []string{"gravl", "strava", "activities", "--filter", ".Type"},
			Err:  "bool",
		},
		{
			Name: "attribute fails on an activity",
			Args: []string{"gravl", "strava", "activities", "--attribute", ".Laps[5].Name"},
			Err:  "index out of range",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}

func TestConcurrency(t *testing.T) {
	a := assert.New(t)
	tests := []*internal.Harness{
		{
			Name:     "zero concurrency still fetches",
			Args:     []string{"gravl", "strava", "--concurrency", "0", "activity", "12345", "54321"},
			Counters: map[string]int{"gravl.strava.activity": 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, activityHandler(a), command)
		})
	}
}

func TestWebhookUnsubscribeNone(t *testing.T) {
	a := assert.New(t)
	mux := http.NewServeMux()
	mux.HandleFunc("/push_subscriptions", func(w http.ResponseWriter, _ *http.Request) {
		a.NoError(json.NewEncoder(w).Encode([]*api.WebhookSubscription{}))
	})
	tt := &internal.Harness{Name: "no subscriptions", Args: []string{"gravl", "strava", "webhook", "unsubscribe"}}
	internal.Run(t, tt, mux, command)
}

// expiredCommand is command with an expired token, so refreshing calls the token endpoint.
func expiredCommand(t *testing.T, baseURL string) *cli.Command {
	endpoint := api.Endpoint()
	endpoint.TokenURL = baseURL + "/token"
	c := strava.Command()
	c.Before = func(c *cli.Context) error {
		client, err := api.NewClient(
			api.WithBaseURL(baseURL),
			api.WithConfig(oauth2.Config{Endpoint: endpoint}),
			api.WithClientCredentials("id", "secret"),
			api.WithTokenCredentials("foo", "bar", time.Now().Add(-time.Hour)))
		if err != nil {
			t.Error(err)
		}
		gravl.Runtime(c).Strava = client
		return nil
	}
	return c
}

func TestRefreshFailure(t *testing.T) {
	tt := &internal.Harness{Name: "refresh", Args: []string{"gravl", "strava", "refresh"}, Err: "oauth2"}
	internal.Run(t, tt, internal.FailingAPI(t), expiredCommand)
}
