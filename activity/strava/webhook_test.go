package strava_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/bzimmer/activity/strava"
	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/internal"
)

func TestWebhookList(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/push_subscriptions", func(w http.ResponseWriter, r *http.Request) {
		res := []*strava.WebhookSubscription{}
		q := r.URL.Query()
		if q.Get("client_id") == "two-subs" {
			res = append(res, []*strava.WebhookSubscription{{}, {}}...)
		}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(res))
	})

	tests := []*internal.Harness{
		{
			Name: "no active subscriptions",
			Args: []string{"gravl", "strava", "webhook", "list"},
		},
		{
			Name: "two active subscriptions",
			Args: []string{"gravl", "strava", "--strava-client-id", "two-subs", "webhook", "list"},
			Counters: map[string]int{
				"gravl.strava.webhook.list": 2,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}

func TestWebhookUnsubscribe(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/push_subscriptions/334455", func(_ http.ResponseWriter, r *http.Request) {
		a.Equal(http.MethodDelete, r.Method)
		// no response
	})
	mux.HandleFunc("/push_subscriptions/10", func(_ http.ResponseWriter, r *http.Request) {
		a.Equal(http.MethodDelete, r.Method)
		// no response
	})
	mux.HandleFunc("/push_subscriptions/20", func(_ http.ResponseWriter, r *http.Request) {
		a.Equal(http.MethodDelete, r.Method)
		// no response
	})
	mux.HandleFunc("/push_subscriptions", func(w http.ResponseWriter, _ *http.Request) {
		res := []*strava.WebhookSubscription{{ID: 10}, {ID: 20}}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(res))
	})

	tests := []*internal.Harness{
		{
			Name: "unsubscribe with no args",
			Args: []string{"gravl", "strava", "webhook", "unsubscribe"},
			Counters: map[string]int{
				"gravl.strava.unsubscribe": 2,
			},
		},
		{
			Name: "unsubscribe one",
			Args: []string{"gravl", "strava", "webhook", "unsubscribe", "334455"},
			Counters: map[string]int{
				"gravl.strava.unsubscribe": 1,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}

func TestWebhookSubscribe(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/push_subscriptions", func(w http.ResponseWriter, r *http.Request) {
		a.Equal(http.MethodPost, r.Method)
		a.NoError(r.ParseForm())
		a.Equal("https://example.com", r.Form.Get("callback_url"))

		res := &strava.WebhookAcknowledgement{ID: 10}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(res))
	})

	tests := []*internal.Harness{
		{
			Name: "subscribe",
			Args: []string{"gravl", "strava", "webhook", "subscribe", "--url", "https://example.com"},
			Counters: map[string]int{
				"gravl.strava.webhook.subscribe": 1,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}
