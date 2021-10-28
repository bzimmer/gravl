package inreach_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"

	api "github.com/bzimmer/activity/inreach"
	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/activity/inreach"
	"github.com/bzimmer/gravl/pkg/internal"
)

func command(t *testing.T, baseURL string) *cli.Command {
	c := inreach.Command()
	c.Before = func(c *cli.Context) error {
		client, err := api.NewClient(
			api.WithBaseURL(baseURL),
			api.WithHTTPTracing(c.Bool("http-tracing")))
		if err != nil {
			t.Error(err)
		}
		pkg.Runtime(c).InReach = client
		return nil
	}
	return c
}

func TestActivity(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/Feed/Share/foobar", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/feed.kml")
	})
	mux.HandleFunc("/Feed/Share/datetimer", func(w http.ResponseWriter, r *http.Request) {
		d1, ok := r.URL.Query()["d1"]
		a.True(ok)
		a.NotNil(d1)
		a.Len(d1, 1)
		d1time, err := time.Parse(api.DateFormat, d1[0])
		a.NoError(err)
		d2, ok := r.URL.Query()["d2"]
		a.True(ok)
		a.NotNil(d2)
		a.Len(d2, 1)
		d2time, err := time.Parse(api.DateFormat, d2[0])
		a.NoError(err)
		a.True(d1time.Before(d2time))
		http.ServeFile(w, r, "testdata/feed.kml")
	})

	tests := []*internal.Harness{
		{
			Name:     "feed",
			Args:     []string{"gravl", "inreach", "feed", "foobar"},
			Counters: map[string]int{"gravl.inreach.feed": 1},
		},
		{
			Name:     "feed since",
			Args:     []string{"gravl", "inreach", "feed", "--since", "72 hours ago", "datetimer"},
			Counters: map[string]int{"gravl.inreach.feed": 1},
		},
		{
			Name: "feed since invalid",
			Args: []string{"gravl", "inreach", "feed", "--since", "72h", "datetimer"},
			Err:  "invalid date range",
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
		Before: inreach.Before,
		Metadata: map[string]interface{}{
			pkg.RuntimeKey: &pkg.Rt{},
		},
		Action: func(c *cli.Context) error {
			a.NotNil(pkg.Runtime(c).InReach)
			return nil
		},
		Commands: []*cli.Command{
			inreach.Command(),
		},
	}
	a.NoError(app.RunContext(context.Background(), []string{"test"}))
}
