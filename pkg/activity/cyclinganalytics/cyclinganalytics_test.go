package cyclinganalytics_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"

	api "github.com/bzimmer/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/pkg/internal"
)

func command(t *testing.T, baseURL string) *cli.Command {
	endpoint := api.Endpoint()
	endpoint.AuthURL = baseURL
	endpoint.TokenURL = baseURL
	c := cyclinganalytics.Command()
	c.Before = func(c *cli.Context) error {
		client, err := api.NewClient(
			api.WithBaseURL(baseURL),
			api.WithHTTPTracing(false),
			api.WithConfig(oauth2.Config{Endpoint: endpoint}),
			api.WithTokenCredentials("foo", "bar", time.Now().Add(time.Hour*24)))
		if err != nil {
			t.Error(err)
		}
		pkg.Runtime(c).CyclingAnalytics = client
		return nil
	}
	return c
}

func TestAthlete(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		ath := &api.User{}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(ath))
	})

	tests := []*internal.Harness{
		{
			Name:     "athlete",
			Args:     []string{"gravl", "cyclinganalytics", "athlete"},
			Counters: map[string]int{"gravl.cyclinganalytics.athlete": 1},
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
			Args:     []string{"gravl", "cyclinganalytics", "streamsets"},
			Counters: map[string]int{"gravl.cyclinganalytics.streamsets": 1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, nil, command)
		})
	}
}

func TestActivity(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/ride/77282721", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(api.Ride{
			LocalDatetime: api.Datetime{Time: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local)},
			UTCDatetime:   api.Datetime{Time: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
		}))
	})

	tests := []*internal.Harness{
		{
			Name:     "ride",
			Args:     []string{"gravl", "cyclinganalytics", "activity", "77282721"},
			Counters: map[string]int{"gravl.cyclinganalytics.activity": 1},
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
			Before: cyclinganalytics.Before,
			Counters: map[string]int{
				"gravl.cyclinganalytics.client.created": 1,
			},
			Action: func(c *cli.Context) error {
				a.NotNil(pkg.Runtime(c).CyclingAnalytics)
				a.NotNil(pkg.Runtime(c).Endpoints[cyclinganalytics.Provider])
				return nil
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			cmd := func(t *testing.T, baseURL string) *cli.Command {
				return &cli.Command{Name: tt.Name, Flags: cyclinganalytics.AuthFlags(), Action: tt.Action}
			}
			internal.Run(t, tt, nil, cmd)
		})
	}
}

func TestActivities(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/me/rides", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(struct {
			Rides []*api.Ride `json:"rides"`
		}{
			Rides: []*api.Ride{{ID: 1008}, {ID: 8001}, {ID: 1772}},
		}))
	})

	tests := []*internal.Harness{
		{
			Name: "ride",
			Args: []string{"gravl", "cyclinganalytics", "activities"},
			Counters: map[string]int{
				"gravl.cyclinganalytics.activity":   3,
				"gravl.cyclinganalytics.activities": 1,
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
