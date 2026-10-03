package rwgps_test

import (
	"encoding/json"
	"net/http"
	"testing"

	api "github.com/bzimmer/activity/rwgps"
	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/internal"
)

func TestAPIFailures(t *testing.T) {
	for _, args := range [][]string{
		{"gravl", "rwgps", "activities"},
		{"gravl", "rwgps", "routes"},
		{"gravl", "rwgps", "activity", "7728201"},
		{"gravl", "rwgps", "route", "90288724"},
	} {
		t.Run(args[2], func(t *testing.T) {
			tt := &internal.Harness{Name: args[2], Args: args, AnyErr: true}
			internal.Run(t, tt, internal.FailingAPI(t), command)
		})
	}
}

func TestListFailures(t *testing.T) {
	a := assert.New(t)
	// the user lookup succeeds, the listing fails
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/users/current.json" {
			a.NoError(json.NewEncoder(w).Encode(struct {
				User *api.User `json:"user"`
			}{User: &api.User{ID: 82877292}}))
			return
		}
		internal.FailingAPI(t).ServeHTTP(w, r)
	})
	for _, args := range [][]string{
		{"gravl", "rwgps", "activities"},
		{"gravl", "rwgps", "routes"},
	} {
		t.Run(args[2], func(t *testing.T) {
			tt := &internal.Harness{Name: args[2], Args: args, AnyErr: true}
			internal.Run(t, tt, handler, command)
		})
	}
}

func TestEncodingFailures(t *testing.T) {
	a := assert.New(t)
	mux := http.NewServeMux()
	mux.HandleFunc("/users/current.json", func(w http.ResponseWriter, _ *http.Request) {
		a.NoError(json.NewEncoder(w).Encode(struct {
			User *api.User `json:"user"`
		}{User: &api.User{ID: 82877292}}))
	})
	results := func(w http.ResponseWriter, _ *http.Request) {
		a.NoError(json.NewEncoder(w).Encode(struct {
			Results      []*api.Trip `json:"results"`
			ResultsCount int         `json:"results_count"`
		}{Results: []*api.Trip{{ID: 90399289}}, ResultsCount: 1}))
	}
	mux.HandleFunc("/users/82877292/trips.json", results)
	mux.HandleFunc("/users/82877292/routes.json", results)
	mux.HandleFunc("/trips/7728201.json", func(w http.ResponseWriter, _ *http.Request) {
		a.NoError(json.NewEncoder(w).Encode(struct {
			Type string    `json:"type"`
			Trip *api.Trip `json:"trip"`
		}{Type: "trip", Trip: &api.Trip{ID: 7728201}}))
	})
	for _, args := range [][]string{
		{"gravl", "rwgps", "athlete"},
		{"gravl", "rwgps", "activities", "-N", "1"},
		{"gravl", "rwgps", "routes", "-N", "1"},
		{"gravl", "rwgps", "activity", "7728201"},
	} {
		t.Run(args[2], func(t *testing.T) {
			tt := &internal.Harness{Name: args[2], Args: args, Before: internal.FailEncoding, Err: internal.ErrEncode.Error()}
			internal.Run(t, tt, mux, command)
		})
	}
}
