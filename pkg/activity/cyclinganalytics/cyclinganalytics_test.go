package cyclinganalytics_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	api "github.com/bzimmer/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/pkg/internal"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

func command(t *testing.T, baseURL string) *cli.Command {
	endpoint := api.Endpoint()
	endpoint.AuthURL = baseURL
	endpoint.TokenURL = baseURL
	c := cyclinganalytics.Command()
	c.Before = func(c *cli.Context) error {
		client, err := api.NewClient(
			api.WithBaseURL(baseURL),
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
