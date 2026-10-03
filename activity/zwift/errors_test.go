package zwift_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"testing"

	api "github.com/bzimmer/activity/zwift"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/internal"
)

func TestAPIFailures(t *testing.T) {
	for _, args := range [][]string{
		{"gravl", "zwift", "athlete"},
		{"gravl", "zwift", "activities"},
		{"gravl", "zwift", "activity", "9001"},
	} {
		t.Run(args[2], func(t *testing.T) {
			tt := &internal.Harness{Name: args[2], Args: args, AnyErr: true}
			internal.Run(t, tt, internal.FailingAPI(t), command)
		})
	}
}

func TestLookupSucceedsFetchFails(t *testing.T) {
	a := assert.New(t)
	// the profile lookup succeeds, the activity requests fail
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/profiles/me" {
			a.NoError(json.NewEncoder(w).Encode(&api.Profile{ID: 101}))
			return
		}
		internal.FailingAPI(t).ServeHTTP(w, r)
	})
	tests := []*internal.Harness{
		{Name: "activities", Args: []string{"gravl", "zwift", "activities"}, AnyErr: true},
		{Name: "activity", Args: []string{"gravl", "zwift", "activity", "9001"}, AnyErr: true},
		{Name: "invalid id", Args: []string{"gravl", "zwift", "activity", "not-a-number"}, Err: "invalid syntax"},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, handler, command)
		})
	}
}

func TestEncodingFailures(t *testing.T) {
	a := assert.New(t)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/profiles/me", func(w http.ResponseWriter, _ *http.Request) {
		a.NoError(json.NewEncoder(w).Encode(&api.Profile{ID: 101}))
	})
	mux.HandleFunc("/api/profiles/101/activities/9001", func(w http.ResponseWriter, _ *http.Request) {
		a.NoError(json.NewEncoder(w).Encode(&api.Activity{ID: 9001}))
	})
	mux.HandleFunc("/api/profiles/101/activities/", func(w http.ResponseWriter, _ *http.Request) {
		a.NoError(json.NewEncoder(w).Encode([]*api.Activity{{ID: 9001}}))
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		token := `{"access_token":"abc123","token_type":"bearer","expires_in":3600,"refresh_token":"def456"}`
		_, err := w.Write([]byte(token))
		a.NoError(err)
	})
	big := func(c *cli.Context) error {
		if err := internal.FailEncoding(c); err != nil {
			return err
		}
		return afero.WriteFile(gravl.Runtime(c).Fs, "/foo/fit/ride.fit", make([]byte, 4096), 0o644)
	}
	tests := []*internal.Harness{
		{Name: "athlete", Args: []string{"gravl", "zwift", "athlete"}},
		{Name: "activities", Args: []string{"gravl", "zwift", "activities", "-N", "1"}},
		{Name: "activity", Args: []string{"gravl", "zwift", "activity", "9001"}},
		{Name: "refresh", Args: []string{"gravl", "zwift", "refresh"}},
		{Name: "files", Args: []string{"gravl", "zwift", "files", "/foo/fit"}, Before: big},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Before == nil {
				tt.Before = internal.FailEncoding
			}
			tt.Err = internal.ErrEncode.Error()
			internal.Run(t, tt, mux, command)
		})
	}
}

var errStat = errors.New("stat failure")

// statFailingFs fails every Stat, as a permission error would.
type statFailingFs struct{ afero.Fs }

func (statFailingFs) Stat(string) (os.FileInfo, error) { return nil, errStat }

func TestFilesFailures(t *testing.T) {
	tests := []*internal.Harness{
		{
			Name: "walk fails",
			Args: []string{"gravl", "zwift", "files", "/foo"},
			Before: func(c *cli.Context) error {
				gravl.Runtime(c).Fs = statFailingFs{gravl.Runtime(c).Fs}
				return nil
			},
			Err: errStat.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, internal.FailingAPI(t), command)
		})
	}
}

func TestFilesWithoutHome(t *testing.T) {
	// with no home directory there is no default Zwift folder, which is not an error
	t.Setenv("HOME", "")
	tt := &internal.Harness{Name: "no home", Args: []string{"gravl", "zwift", "files"}}
	internal.Run(t, tt, internal.FailingAPI(t), command)
}

func TestRefreshFailure(t *testing.T) {
	tt := &internal.Harness{Name: "refresh", Args: []string{"gravl", "zwift", "refresh"}, AnyErr: true}
	internal.Run(t, tt, internal.FailingAPI(t), command)
}
