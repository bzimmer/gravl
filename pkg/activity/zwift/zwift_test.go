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
)

func command(t *testing.T, baseURL string) *cli.Command {
	c := zwift.Command()
	c.Before = func(c *cli.Context) error {
		client, err := api.NewClient(
			api.WithBaseURL(baseURL),
			api.WithTokenCredentials("foo", "bar", time.Now()))
		if err != nil {
			t.Error(err)
		}
		pkg.Runtime(c).Zwift = client
		return nil
	}
	return c
}

func TestAthlete(t *testing.T) {
	t.Parallel()
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
