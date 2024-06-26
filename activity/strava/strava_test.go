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

func command(t *testing.T, baseURL string) *cli.Command {
	endpoint := api.Endpoint()
	endpoint.AuthURL = baseURL + "/auth"
	endpoint.TokenURL = baseURL + "/token"
	c := strava.Command()
	c.Before = func(c *cli.Context) error {
		client, err := api.NewClient(
			api.WithBaseURL(baseURL),
			api.WithHTTPTracing(c.Bool("http-tracing")),
			api.WithConfig(oauth2.Config{Endpoint: endpoint}),
			api.WithClientCredentials(c.String("strava-client-id"), "dummy"),
			api.WithTokenCredentials("foo", "bar", time.Now().Add(time.Hour*24)))
		if err != nil {
			t.Error(err)
		}
		gravl.Runtime(c).Strava = client
		return nil
	}
	return c
}

func TestAthlete(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/athlete", func(w http.ResponseWriter, _ *http.Request) {
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
	mux.HandleFunc("/activities/12345", func(w http.ResponseWriter, _ *http.Request) {
		act := &api.Activity{}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(act))
	})
	mux.HandleFunc("/activities/12345/streams/", func(w http.ResponseWriter, _ *http.Request) {
		sms := &api.Streams{}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(sms))
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
		{
			Name: "no arguments",
			Args: []string{"gravl", "strava", "activity"},
		},
		{
			Name: "activity streams",
			Args: []string{"gravl", "strava", "activity", "-s", "latlng", "12345"},
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
	mux.HandleFunc("/athlete/activities", func(w http.ResponseWriter, _ *http.Request) {
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
			Name: "activities since invalid",
			Args: []string{"gravl", "strava", "activities", "-N", "3", "--since", "72h"},
			Err:  "invalid date range",
		},
		{
			Name:     "activities since",
			Args:     []string{"gravl", "strava", "activities", "-N", "3", "--since", "2 weeks ago"},
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
			Name:     "streamsets",
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
	mux.HandleFunc("/routes/77282721", func(w http.ResponseWriter, _ *http.Request) {
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
	mux.HandleFunc("/activities/77282721/streams/", func(w http.ResponseWriter, _ *http.Request) {
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
	mux.HandleFunc("/activities/26627201/photos", func(w http.ResponseWriter, _ *http.Request) {
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
	mux.HandleFunc("/athlete", func(w http.ResponseWriter, _ *http.Request) {
		ath := &api.Athlete{ID: 8542982}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(ath))
	})
	mux.HandleFunc("/athletes/8542982/routes", func(w http.ResponseWriter, _ *http.Request) {
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
	tests := []*internal.Harness{
		{
			Name:   "testbefore",
			Args:   []string{"gravl", "testbefore"},
			Before: gravl.Befores(strava.Before, strava.Before, strava.Before, strava.Before),
			Counters: map[string]int{
				"gravl.strava.client.created": 1,
			},
			Action: func(c *cli.Context) error {
				a.NotNil(gravl.Runtime(c).Strava)
				a.NotNil(gravl.Runtime(c).Endpoints[strava.Provider])
				return nil
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			cmd := func(_ *testing.T, _ string) *cli.Command {
				return &cli.Command{Name: tt.Name, Flags: strava.AuthFlags(), Action: tt.Action}
			}
			internal.Run(t, tt, nil, cmd)
		})
	}
}

func TestRefresh(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, _ *http.Request) {
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

func TestUpdate(t *testing.T) {
	a := assert.New(t)

	var decoder = func(r *http.Request) *api.UpdatableActivity {
		var act api.UpdatableActivity
		dec := json.NewDecoder(r.Body)
		a.NoError(dec.Decode(&act))
		return &act
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/activities/101", func(_ http.ResponseWriter, r *http.Request) {
		a.Equal(http.MethodPut, r.Method)
		act := decoder(r)
		a.True(*act.Hidden)
	})
	mux.HandleFunc("/activities/102", func(_ http.ResponseWriter, r *http.Request) {
		a.Equal(http.MethodPut, r.Method)
		act := decoder(r)
		a.False(*act.Hidden)
	})
	mux.HandleFunc("/activities/103", func(_ http.ResponseWriter, r *http.Request) {
		a.Equal(http.MethodPut, r.Method)
		act := decoder(r)
		a.Equal("foobar", *act.Description)
	})
	mux.HandleFunc("/activities/104", func(_ http.ResponseWriter, r *http.Request) {
		a.Equal(http.MethodPut, r.Method)
		act := decoder(r)
		a.Equal("foobaz", *act.Name)
	})
	mux.HandleFunc("/activities/105", func(_ http.ResponseWriter, r *http.Request) {
		a.Equal(http.MethodPut, r.Method)
		act := decoder(r)
		a.Equal("gravel", *act.SportType)
		a.Equal("99181", *act.GearID)
	})
	mux.HandleFunc("/activities/106", func(_ http.ResponseWriter, r *http.Request) {
		a.Equal(http.MethodPut, r.Method)
		act := decoder(r)
		a.True(*act.Commute)
		a.True(*act.Trainer)
	})
	mux.HandleFunc("/activities/107", func(_ http.ResponseWriter, r *http.Request) {
		a.Equal(http.MethodPut, r.Method)
		act := decoder(r)
		a.False(*act.Commute)
		a.False(*act.Trainer)
	})

	tests := []*internal.Harness{
		{
			Name: "update help",
			Args: []string{"gravl", "strava", "update", "--help"},
		},
		{
			Name: "update hidden",
			Args: []string{"gravl", "strava", "update", "--hidden", "101"},
		},
		{
			Name: "update no-hidden",
			Args: []string{"gravl", "strava", "update", "--no-hidden", "102"},
		},
		{
			Name: "update description",
			Args: []string{"gravl", "strava", "update", "--description", "foobar", "103"},
		},
		{
			Name: "update name",
			Args: []string{"gravl", "strava", "update", "--name", "foobaz", "104"},
			Counters: map[string]int{
				"gravl.strava.update":      1,
				"gravl.strava.update.name": 1,
			},
		},
		{
			Name: "update sport & gear",
			Args: []string{"gravl", "strava", "update", "--gear", "99181", "--sport", "gravel", "105"},
		},
		{
			Name: "set trainer and commute",
			Args: []string{"gravl", "strava", "update", "--trainer", "--commute", "106"},
			Counters: map[string]int{
				"gravl.strava.update":         1,
				"gravl.strava.update.trainer": 1,
				"gravl.strava.update.commute": 1,
			},
		},
		{
			Name: "unset trainer and commute",
			Args: []string{"gravl", "strava", "update", "--no-trainer", "--no-commute", "107"},
		},
		{
			Name: "invalid hidden",
			Args: []string{"gravl", "strava", "update", "--hidden", "--no-hidden", "9001"},
			Err:  "only one of hidden or no-hidden can be specified",
		},
		{
			Name: "invalid trainer",
			Args: []string{"gravl", "strava", "update", "--trainer", "--no-trainer", "9001"},
			Err:  "only one of trainer or no-trainer can be specified",
		},
		{
			Name: "invalid commute",
			Args: []string{"gravl", "strava", "update", "--commute", "--no-commute", "9001"},
			Err:  "only one of commute or no-commute can be specified",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}
