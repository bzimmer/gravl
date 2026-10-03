package hammerhead_test

import (
	"encoding/json"
	"net/http"
	"testing"

	api "github.com/bzimmer/activity/hammerhead"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/internal"
)

func readOnly(c *cli.Context) error {
	gravl.Runtime(c).Fs = afero.NewReadOnlyFs(gravl.Runtime(c).Fs)
	return nil
}

func handler(a *assert.Assertions) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, _ *http.Request) {
		a.NoError(json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "newaccesstoken",
			"token_type":    "bearer",
			"expires_in":    3600,
			"refresh_token": "newrefreshtoken",
		}))
	})
	mux.HandleFunc("/activities", func(w http.ResponseWriter, _ *http.Request) {
		a.NoError(json.NewEncoder(w).Encode(&api.ActivitiesPage{
			TotalItems: 1, TotalPages: 1, PerPage: 100, CurrentPage: 1,
			Data: []*api.ActivitySummary{{ID: "act-001", Name: "Morning Ride"}},
		}))
	})
	file := func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.ant.fit")
		_, _ = w.Write([]byte("FIT file content"))
	}
	mux.HandleFunc("/activities/act-001/file", file)
	mux.HandleFunc("/activities/act-002/file", file)
	return mux
}

func TestRefreshCacheFailure(t *testing.T) {
	a := assert.New(t)
	// a token the cache cannot store is still a successful refresh
	tt := &internal.Harness{
		Name:   "read-only cache",
		Args:   []string{"gravl", "hammerhead", "refresh"},
		Before: readOnly,
	}
	internal.Run(t, tt, handler(a), command)
}

func TestEncodingFailures(t *testing.T) {
	a := assert.New(t)
	for _, args := range [][]string{
		{"gravl", "hammerhead", "refresh"},
		{"gravl", "hammerhead", "activities"},
		{"gravl", "hammerhead", "file", "-O", "/tmp/act-001.fit", "act-001"},
	} {
		t.Run(args[2], func(t *testing.T) {
			tt := &internal.Harness{Name: args[2], Args: args, Before: internal.FailEncoding, Err: internal.ErrEncode.Error()}
			internal.Run(t, tt, handler(a), command)
		})
	}
}

func TestFileWriteFailures(t *testing.T) {
	a := assert.New(t)
	tests := []*internal.Harness{
		{
			Name:   "create fails",
			Args:   []string{"gravl", "hammerhead", "file", "-O", "/tmp/act-001.fit", "act-001"},
			Before: readOnly,
			AnyErr: true,
		},
		{
			Name:   "output directory fails",
			Args:   []string{"gravl", "hammerhead", "file", "-O", "/tmp/fits", "act-001", "act-002"},
			Before: readOnly,
			AnyErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, handler(a), command)
		})
	}
}

func TestFileTruncated(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/activities/act-001/file", func(w http.ResponseWriter, _ *http.Request) {
		// promise more than is sent, so reading the body fails partway
		w.Header().Set("Content-Type", "application/vnd.ant.fit")
		w.Header().Set("Content-Length", "1000")
		_, _ = w.Write([]byte("FIT"))
	})
	tt := &internal.Harness{
		Name:   "truncated download",
		Args:   []string{"gravl", "hammerhead", "file", "-O", "/tmp/act-001.fit", "act-001"},
		AnyErr: true,
	}
	internal.Run(t, tt, mux, command)
}
