package strava_test

import (
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
)

func command(t *testing.T, baseURL string) *cli.Command {
	c := strava.Command()
	c.Before = func(c *cli.Context) error {
		client, err := api.NewClient(
			api.WithTokenCredentials("foo", "bar", time.Now()),
			api.WithBaseURL(baseURL))
		if err != nil {
			t.Error(err)
		}
		pkg.Runtime(c).Strava = client
		return nil
	}
	return c
}

func TestAthlete(t *testing.T) {
	t.Parallel()
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
