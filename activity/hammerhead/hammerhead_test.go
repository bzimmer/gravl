package hammerhead_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	api "github.com/bzimmer/activity/hammerhead"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/activity/hammerhead"
	"github.com/bzimmer/gravl/internal"
)

func command(t *testing.T, baseURL string) *cli.Command {
	t.Helper()
	c := hammerhead.Command()
	c.Before = func(c *cli.Context) error {
		client, err := api.NewClient(
			api.WithHTTPTracing(c.Bool("http-tracing")),
			api.WithTokenCredentials("testtoken", "testrefresh", time.Now().Add(time.Hour)),
			api.WithAPIURL(baseURL),
			api.WithAuthURL(baseURL),
		)
		if err != nil {
			t.Error(err)
		}
		gravl.Runtime(c).Hammerhead = client
		return nil
	}
	return c
}

func TestBefore(t *testing.T) {
	a := assert.New(t)
	tests := []*internal.Harness{
		{
			Name:   "testbefore",
			Args:   []string{"gravl", "testbefore"},
			Before: hammerhead.Before,
			Counters: map[string]int{
				"gravl.hammerhead.client.created": 1,
			},
			Action: func(c *cli.Context) error {
				a.NotNil(gravl.Runtime(c).Hammerhead)
				ep, ok := gravl.Runtime(c).Endpoints[hammerhead.Provider]
				a.True(ok)
				a.NotEmpty(ep.AuthURL)
				return nil
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			cmd := func(_ *testing.T, _ string) *cli.Command {
				return &cli.Command{Name: tt.Name, Flags: hammerhead.AuthFlags(), Action: tt.Action}
			}
			internal.Run(t, tt, nil, cmd)
		})
	}
}

func TestActivities(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/activities", func(w http.ResponseWriter, _ *http.Request) {
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(&api.ActivitiesPage{
			TotalItems:  2,
			TotalPages:  1,
			PerPage:     100,
			CurrentPage: 1,
			Data: []*api.ActivitySummary{
				{ID: "act-001", Name: "Morning Ride", Duration: 3600},
				{ID: "act-002", Name: "Evening Gravel", Duration: 5400},
			},
		}))
	})

	tests := []*internal.Harness{
		{
			Name: "activities list",
			Args: []string{"gravl", "hammerhead", "activities"},
			Counters: map[string]int{
				"gravl.hammerhead.activities": 1,
				"gravl.hammerhead.activity":   2,
			},
		},
		{
			Name: "activities with count",
			Args: []string{"gravl", "hammerhead", "activities", "-N", "1"},
			Counters: map[string]int{
				"gravl.hammerhead.activities": 1,
				"gravl.hammerhead.activity":   1,
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

func TestActivity(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/activities/act-001", func(w http.ResponseWriter, _ *http.Request) {
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(&api.Activity{
			ActivitySummary: api.ActivitySummary{
				ID:   "act-001",
				Name: "Morning Ride",
			},
			ActivityType: api.ActivityTypeRide,
		}))
	})

	tests := []*internal.Harness{
		{
			Name:     "activity query",
			Args:     []string{"gravl", "hammerhead", "activity", "act-001"},
			Counters: map[string]int{"gravl.hammerhead.activity": 1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}

func TestExport(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/activities/act-001/file", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.ant.fit")
		_, _ = w.Write([]byte("FIT file content"))
	})
	mux.HandleFunc("/activities/bad/file", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	tests := []*internal.Harness{
		{
			Name:     "export to stdout",
			Args:     []string{"gravl", "hammerhead", "export", "act-001"},
			Counters: map[string]int{"gravl.hammerhead.export": 1},
		},
		{
			Name: "export to file",
			Args: []string{"gravl", "hammerhead", "export", "-O", "/tmp/act-001.fit", "act-001"},
			After: func(c *cli.Context) error {
				stat, err := gravl.Runtime(c).Fs.Stat("/tmp/act-001.fit")
				a.NoError(err)
				a.NotNil(stat)
				return nil
			},
			Counters: map[string]int{"gravl.hammerhead.export": 1},
		},
		{
			Name: "export file exists error",
			Args: []string{"gravl", "hammerhead", "export", "-O", "/tmp/existing.fit", "act-001"},
			Before: func(c *cli.Context) error {
				fp, err := gravl.Runtime(c).Fs.Create("/tmp/existing.fit")
				a.NoError(err)
				return fp.Close()
			},
			Err:      "file already exists",
			Counters: map[string]int{},
		},
		{
			Name: "export with overwrite",
			Args: []string{"gravl", "hammerhead", "export", "-O", "/tmp/overwrite.fit", "-o", "act-001"},
			Before: func(c *cli.Context) error {
				fp, err := gravl.Runtime(c).Fs.Create("/tmp/overwrite.fit")
				a.NoError(err)
				return fp.Close()
			},
			After: func(c *cli.Context) error {
				stat, err := gravl.Runtime(c).Fs.Stat("/tmp/overwrite.fit")
				a.NoError(err)
				a.NotNil(stat)
				return nil
			},
			Counters: map[string]int{"gravl.hammerhead.export": 1},
		},
		{
			Name:     "export error",
			Args:     []string{"gravl", "hammerhead", "export", "bad"},
			Err:      "Not Found",
			Counters: map[string]int{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}

func TestFileMultiple(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/activities/act-001/file", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.ant.fit")
		_, _ = w.Write([]byte("FIT file content 1"))
	})
	mux.HandleFunc("/activities/act-002/file", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.ant.fit")
		_, _ = w.Write([]byte("FIT file content 2"))
	})

	tests := []*internal.Harness{
		{
			Name:     "multiple ids write to disk instead of stdout",
			Args:     []string{"gravl", "hammerhead", "export", "act-001", "act-002"},
			Counters: map[string]int{"gravl.hammerhead.export": 2},
			After: func(c *cli.Context) error {
				fs := gravl.Runtime(c).Fs
				data, err := afero.ReadFile(fs, "act-001.fit")
				a.NoError(err)
				a.Equal("FIT file content 1", string(data))
				data, err = afero.ReadFile(fs, "act-002.fit")
				a.NoError(err)
				a.Equal("FIT file content 2", string(data))
				return nil
			},
		},
		{
			Name:     "output flag rejected with multiple ids",
			Err:      "--output cannot be used with more than one ACTIVITY_ID",
			Args:     []string{"gravl", "hammerhead", "export", "-O", "out.fit", "act-001", "act-002"},
			Counters: map[string]int{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}

func TestRefresh(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, _ *http.Request) {
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(map[string]any{
			"access_token":  "newaccesstoken",
			"token_type":    "bearer",
			"expires_in":    3600,
			"refresh_token": "newrefreshtoken",
		}))
	})

	tests := []*internal.Harness{
		{
			Name: "refresh",
			Args: []string{"gravl", "hammerhead", "refresh"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}

func TestActivityError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/activities/missing", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	tests := []*internal.Harness{
		{
			Name:     "activity not found",
			Err:      "Not Found",
			Args:     []string{"gravl", "hammerhead", "activity", "missing"},
			Counters: map[string]int{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}

func TestActivitiesError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/activities", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	tests := []*internal.Harness{
		{
			Name:     "activities server error",
			Err:      "Internal Server Error",
			Args:     []string{"gravl", "hammerhead", "activities"},
			Counters: map[string]int{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, mux, command)
		})
	}
}

