package rwgps_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	api "github.com/bzimmer/activity/rwgps"
	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/activity/rwgps"
	"github.com/bzimmer/gravl/pkg/internal"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

func command(t *testing.T, baseURL string) *cli.Command {
	c := rwgps.Command()
	c.Before = func(c *cli.Context) error {
		client, err := api.NewClient(
			api.WithHTTPTracing(c.Bool("http-tracing")),
			api.WithTokenCredentials("foo", "bar", time.Now()),
			api.WithBaseURL(baseURL))
		if err != nil {
			t.Error(err)
		}
		pkg.Runtime(c).RideWithGPS = client
		return nil
	}
	return c
}

func TestAthlete(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/users/current.json", func(w http.ResponseWriter, r *http.Request) {
		type res struct {
			User *api.User `json:"user"`
		}
		ath := &api.User{ID: 100, Name: "foo"}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(&res{User: ath}))
	})

	tests := []*internal.Harness{
		{
			Name:     "athlete",
			Args:     []string{"gravl", "rwgps", "athlete"},
			Counters: map[string]int{"gravl.rwgps.athlete": 1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}

func TestBefore(t *testing.T) {
	a := assert.New(t)
	app := &cli.App{
		Name:   "TestBefore",
		Before: rwgps.Before,
		Metadata: map[string]interface{}{
			pkg.RuntimeKey: &pkg.Rt{
				Endpoints: make(map[string]oauth2.Endpoint),
			},
		},
		Action: func(c *cli.Context) error {
			a.NotNil(pkg.Runtime(c).RideWithGPS)
			return nil
		},
	}
	a.NoError(app.RunContext(context.Background(), []string{"test"}))
}

func TestTrip(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/routes/90288724.json", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(struct {
			Type  string    `json:"type"`
			Route *api.Trip `json:"route"`
		}{
			Type: "route",
			Route: &api.Trip{
				ID: 90288724,
			},
		}))
	})
	mux.HandleFunc("/trips/7728201.json", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(struct {
			Type string    `json:"type"`
			Trip *api.Trip `json:"trip"`
		}{
			Type: "trip",
			Trip: &api.Trip{
				ID: 7728201,
			},
		}))
	})

	tests := []*internal.Harness{
		{
			Name:     "activity query",
			Args:     []string{"gravl", "rwgps", "activity", "7728201"},
			Counters: map[string]int{"gravl.rwgps.activity": 1},
		},
		{
			Name:     "route query",
			Args:     []string{"gravl", "rwgps", "route", "90288724"},
			Counters: map[string]int{"gravl.rwgps.route": 1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}

func TestTrips(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/users/current.json", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(struct {
			User *api.User `json:"user"`
		}{
			User: &api.User{ID: 82877292},
		}))
	})
	mux.HandleFunc("/users/82877292/trips.json", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(struct {
			Results      []*api.Trip `json:"results"`
			ResultsCount int         `json:"results_count"`
		}{
			Results: []*api.Trip{
				{ID: 90399289},
				{ID: 82827929},
			},
			ResultsCount: 2,
		}))
	})
	mux.HandleFunc("/users/82877292/routes.json", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(struct {
			Results      []*api.Trip `json:"results"`
			ResultsCount int         `json:"results_count"`
		}{
			Results: []*api.Trip{
				{ID: 22233292},
				{ID: 82823222},
				{ID: 28839283},
				{ID: 53202008},
			},
			ResultsCount: 4,
		}))
	})

	tests := []*internal.Harness{
		{
			Name: "activities two",
			Args: []string{"gravl", "rwgps", "activities", "-N", "2"},
			Counters: map[string]int{
				"gravl.rwgps.activities": 1,
				"gravl.rwgps.activity":   2,
			},
		},
		{
			Name: "routes two",
			Args: []string{"gravl", "rwgps", "routes", "-N", "2"},
			Counters: map[string]int{
				"gravl.rwgps.routes": 1,
				"gravl.rwgps.route":  2,
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
