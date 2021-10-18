package strava_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	api "github.com/bzimmer/activity/strava"
	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/activity/strava"
	"github.com/bzimmer/gravl/pkg/internal"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

func command(t *testing.T, baseURL string) *cli.Command {
	endpoint := api.Endpoint()
	endpoint.AuthURL = baseURL + "/auth"
	endpoint.TokenURL = baseURL + "/token"
	c := strava.Command()
	c.Before = func(c *cli.Context) error {
		client, err := api.NewClient(
			api.WithBaseURL(baseURL),
			api.WithHTTPTracing(true),
			api.WithConfig(oauth2.Config{Endpoint: endpoint}),
			api.WithTokenCredentials("foo", "bar", time.Now().Add(time.Hour*24)))
		if err != nil {
			t.Error(err)
		}
		pkg.Runtime(c).Strava = client
		return nil
	}
	return c
}

func TestAthlete(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/athlete", func(w http.ResponseWriter, r *http.Request) {
		ath := &api.Athlete{}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(ath))
	})

	tests := []*internal.Harness{
		{
			Name:     "athlete",
			Args:     []string{"gravl", "strava", "athlete"},
			Counters: map[string]int{"gravl.strava.athlete": 1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}

func TestActivity(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/activities/12345", func(w http.ResponseWriter, r *http.Request) {
		act := &api.Activity{}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(act))
	})

	tests := []*internal.Harness{
		{
			Name:     "activity",
			Args:     []string{"gravl", "strava", "activity", "12345"},
			Counters: map[string]int{"gravl.strava.activity": 1},
		},
		{
			Name: "invalid syntax",
			Args: []string{"gravl", "strava", "activity", "12345", "54321", "abcdef"},
			Err:  "invalid syntax",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}

func TestActivities(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/athlete/activities", func(w http.ResponseWriter, r *http.Request) {
		acts := []*api.Activity{{Type: "Hike"}, {Type: "Run"}, {Type: "Hike"}}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(acts))
	})

	tests := []*internal.Harness{
		{
			Name:     "activities",
			Args:     []string{"gravl", "strava", "activities", "-N", "2"},
			Counters: map[string]int{"gravl.strava.activity": 2},
		},
		{
			Name:     "filtered activities",
			Args:     []string{"gravl", "strava", "activities", "-N", "3", "--filter", ".Type == 'Hike'"},
			Counters: map[string]int{"gravl.strava.activity": 2},
		},
		{
			Name:     "activity attributes",
			Args:     []string{"gravl", "strava", "activities", "-N", "3", "--attribute", ".Type"},
			Counters: map[string]int{"gravl.strava.activity": 3},
		},
		{
			Name:     "activities since",
			Args:     []string{"gravl", "strava", "activities", "-N", "3", "--since", "72h"},
			Counters: map[string]int{"gravl.strava.activity": 3},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}

func TestStreamSets(t *testing.T) {
	tests := []*internal.Harness{
		{
			Name:     "activity",
			Args:     []string{"gravl", "strava", "streamsets"},
			Counters: map[string]int{"gravl.strava.streamsets": 1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, nil, command)
		})
	}
}

func TestRoute(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/routes/77282721", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(api.Route{}))
	})

	tests := []*internal.Harness{
		{
			Name:     "routes",
			Args:     []string{"gravl", "strava", "route", "77282721"},
			Counters: map[string]int{"gravl.strava.route": 1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}

func TestStreams(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/activities/77282721/streams/", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(api.Streams{}))
	})

	tests := []*internal.Harness{
		{
			Name: "invalid streams",
			Args: []string{"gravl", "strava", "streams", "-s", "latlng", "-s", "foobar", "77282721"},
			Err:  "invalid stream",
		},
		{
			Name:     "streams",
			Args:     []string{"gravl", "strava", "streams", "-s", "latlng", "-s", "temp", "77282721"},
			Counters: map[string]int{"gravl.strava.streams": 1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}

func TestPhotos(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/activities/26627201/photos", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode([]*api.Photo{{}, {}}))
	})

	tests := []*internal.Harness{
		{
			Name:     "photos",
			Args:     []string{"gravl", "strava", "photos", "26627201"},
			Counters: map[string]int{"gravl.strava.photos": 1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}

func TestRoutes(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/athlete", func(w http.ResponseWriter, r *http.Request) {
		ath := &api.Athlete{ID: 8542982}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(ath))
	})
	mux.HandleFunc("/athletes/8542982/routes", func(w http.ResponseWriter, r *http.Request) {
		rtes := []*api.Route{{}, {}, {}}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(rtes))
	})

	tests := []*internal.Harness{
		{
			Name: "routes",
			Args: []string{"gravl", "strava", "routes", "-N", "2"},
			Counters: map[string]int{
				"gravl.strava.route":  2,
				"gravl.strava.routes": 1,
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

func TestBefore(t *testing.T) {
	a := assert.New(t)
	app := &cli.App{
		Name:   "TestBefore",
		Before: strava.Before,
		Metadata: map[string]interface{}{
			pkg.RuntimeKey: &pkg.Rt{},
		},
		Action: func(c *cli.Context) error {
			a.NotNil(pkg.Runtime(c).Strava)
			return nil
		},
	}
	a.NoError(app.RunContext(context.Background(), []string{"test"}))
}

func TestRefresh(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		n, err := w.Write([]byte(`{
			"access_token":"90d64460d14870c08c81352a05dedd3465940a7c",
			"token_type":"bearer",
			"expires_in":3600,
			"refresh_token":"IwOGYzYTlmM2YxOTQ5MGE3YmNmMDFkNTVk",
			"scope":"user"
		  }`))
		a.Greater(n, 0)
		a.NoError(err)
	})

	tests := []*internal.Harness{
		{
			Name: "refresh",
			Args: []string{"gravl", "strava", "refresh"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}
