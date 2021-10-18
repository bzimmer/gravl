package zwift_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	api "github.com/bzimmer/activity/zwift"
	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/activity/zwift"
	"github.com/bzimmer/gravl/pkg/internal"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

func command(t *testing.T, baseURL string) *cli.Command {
	endpoint := api.Endpoint()
	endpoint.AuthURL = baseURL + "/auth"
	endpoint.TokenURL = baseURL + "/token"
	c := zwift.Command()
	c.Before = func(c *cli.Context) error {
		client, err := api.NewClient(
			api.WithBaseURL(baseURL),
			api.WithHTTPTracing(false),
			api.WithConfig(oauth2.Config{Endpoint: endpoint}),
			api.WithTokenCredentials("foo", "bar", time.Now().Add(time.Hour*24)))
		if err != nil {
			t.Error(err)
		}
		pkg.Runtime(c).Zwift = client
		return nil
	}
	return c
}

func TestAthlete(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/profiles/me", func(w http.ResponseWriter, r *http.Request) {
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
				fs := pkg.Runtime(c).Fs
				a.NoError(fs.MkdirAll("/foo/bar/Zwift/Activities", 0777))
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
