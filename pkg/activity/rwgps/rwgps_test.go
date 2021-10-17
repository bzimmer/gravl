package rwgps_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	api "github.com/bzimmer/activity/rwgps"
	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/activity/rwgps"
	"github.com/bzimmer/gravl/pkg/internal"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func command(t *testing.T, baseURL string) *cli.Command {
	c := rwgps.Command()
	c.Before = func(c *cli.Context) error {
		client, err := api.NewClient(
			api.WithTokenCredentials("foo", "bar", time.Now()),
			api.WithBaseURL(baseURL))
		if err != nil {
			t.Error(err)
		}
		pkg.Runtime(c).RideWithGPS = client
		return nil
	}
	return c
}

func TestAthlete(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/users/current.json", func(w http.ResponseWriter, r *http.Request) {
		ath := &api.User{}
		enc := json.NewEncoder(w)
		a.NoError(enc.Encode(ath))
	})

	tests := []*internal.Harness{
		{
			Name:     "athlete",
			Args:     []string{"gravl", "rwgps", "athlete"},
			Counters: map[string]int{"gravl.rwgps.athlete": 1},
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
	app := &cli.App{
		Name:   "TestBefore",
		Before: rwgps.Before,
		Metadata: map[string]interface{}{
			pkg.RuntimeKey: &pkg.Rt{},
		},
		Action: func(c *cli.Context) error {
			a.NotNil(pkg.Runtime(c).RideWithGPS)
			return nil
		},
	}
	a.NoError(app.RunContext(context.Background(), []string{"test"}))
}
