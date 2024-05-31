package zwift_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	api "github.com/bzimmer/activity/zwift"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/activity/zwift"
	"github.com/bzimmer/gravl/internal"
)

func command(t *testing.T, baseURL string) *cli.Command {
	endpoint := api.Endpoint()
	endpoint.AuthURL = baseURL + "/auth"
	endpoint.TokenURL = baseURL + "/token"
	c := zwift.Command()
	c.Before = func(c *cli.Context) error {
		client, err := api.NewClient(
			api.WithBaseURL(baseURL),
			api.WithHTTPTracing(c.Bool("http-tracing")),
			api.WithConfig(oauth2.Config{Endpoint: endpoint}),
			api.WithTokenCredentials("foo", "bar", time.Now().Add(time.Hour*24)))
		if err != nil {
			t.Error(err)
		}
		gravl.Runtime(c).Zwift = client
		return nil
	}
	return c
}

func TestAthlete(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/profiles/me", func(w http.ResponseWriter, _ *http.Request) {
		ath := &api.Profile{}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(ath))
	})

	tests := []*internal.Harness{
		{
			Name:     "me",
			Args:     []string{"gravl", "zwift", "athlete"},
			Counters: map[string]int{"gravl.zwift.athlete": 1},
		},
		{
			Name: "unknown athlete",
			Err:  http.StatusText(http.StatusNotFound),
			Args: []string{"gravl", "zwift", "athlete", "foobar"},
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
	mux.HandleFunc("/api/profiles/me", func(w http.ResponseWriter, _ *http.Request) {
		ath := &api.Profile{ID: 101}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(ath))
	})
	mux.HandleFunc("/api/profiles/101/activities/9001", func(w http.ResponseWriter, _ *http.Request) {
		act := &api.Activity{ID: 9001}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(act))
	})

	tests := []*internal.Harness{
		{
			Name:     "single activity",
			Args:     []string{"gravl", "zwift", "activity", "9001"},
			Counters: map[string]int{"gravl.zwift.activity": 1},
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
	mux.HandleFunc("/api/profiles/me", func(w http.ResponseWriter, _ *http.Request) {
		ath := &api.Profile{ID: 102}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(ath))
	})
	mux.HandleFunc("/api/profiles/102/activities/", func(w http.ResponseWriter, _ *http.Request) {
		acts := []*api.Activity{{ID: 9001}, {ID: 9021}, {ID: 9501}}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(acts))
	})

	tests := []*internal.Harness{
		{
			Name: "two activities",
			Args: []string{"gravl", "zwift", "activities", "-N", "2"},
			Counters: map[string]int{
				"gravl.zwift.activity":   2,
				"gravl.zwift.activities": 1,
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

func TestFiles(t *testing.T) {
	tests := []*internal.Harness{
		{
			Name: "files with specific path",
			Args: []string{"gravl", "zwift", "files", "/foo/bar"},
			Counters: map[string]int{
				"gravl.zwift.files.found":                4,
				"gravl.zwift.files.directory":            3,
				"gravl.zwift.files.skipping.in-progress": 1,
			},
			Before: func(c *cli.Context) error {
				a := assert.New(t)
				fs := gravl.Runtime(c).Fs
				a.NoError(fs.MkdirAll("/foo/bar/Zwift/Activities", 0755))
				fp, err := fs.Create("/foo/bar/Zwift/Activities/inProgressActivity.fit")
				a.NoError(err)
				return fp.Close()
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, nil, command)
		})
	}
}

func TestBefore(t *testing.T) {
	a := assert.New(t)
	tests := []*internal.Harness{
		{
			Name:   "testbefore",
			Args:   []string{"gravl", "testbefore"},
			Before: zwift.Before,
			Counters: map[string]int{
				"gravl.zwift.client.created": 1,
			},
			Action: func(c *cli.Context) error {
				a.NotNil(gravl.Runtime(c).Zwift)
				a.NotNil(gravl.Runtime(c).Endpoints[zwift.Provider])
				return nil
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			cmd := func(_ *testing.T, _ string) *cli.Command {
				return &cli.Command{Name: tt.Name, Flags: zwift.AuthFlags(), Action: tt.Action}
			}
			internal.Run(t, tt, nil, cmd)
		})
	}
}
