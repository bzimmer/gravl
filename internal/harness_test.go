package internal_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/internal"
)

func command(_ *testing.T, _ string) *cli.Command {
	return &cli.Command{
		Name: "foo",
		Before: func(c *cli.Context) error {
			gravl.Runtime(c).Metrics.IncrCounter([]string{c.Command.Name, "before"}, 1)
			return nil
		},
		After: func(c *cli.Context) error {
			gravl.Runtime(c).Metrics.IncrCounter([]string{c.Command.Name, "after"}, 1)
			return nil
		},
		Action: func(c *cli.Context) error {
			gravl.Runtime(c).Metrics.IncrCounter([]string{c.Command.Name, "action"}, 1)
			return nil
		},
	}
}

func TestHarness(t *testing.T) {
	tests := []*internal.Harness{
		{
			Name: "harness",
			Args: []string{"gravl", "foo"},
			Counters: map[string]int{
				"gravl.app.before.tt": 1,
				"gravl.foo.action":    1,
				"gravl.foo.after":     1,
				"gravl.foo.before":    1,
			},
			Before: func(c *cli.Context) error {
				gravl.Runtime(c).Metrics.IncrCounter([]string{"app", "before", "tt"}, 1)
				return nil
			},
			After: func(_ *cli.Context) error {
				return nil
			},
		},
		{
			Name:     "harness with err",
			Args:     []string{"gravl", "--json", "foo"},
			Err:      "foo err bar",
			Counters: map[string]int{},
			Before: func(c *cli.Context) error {
				gravl.Runtime(c).Metrics.IncrCounter([]string{"app", "before", "tt"}, 1)
				return nil
			},
			After: func(_ *cli.Context) error {
				return errors.New("foo err bar")
			},
		},
		{
			Name: "harness no sample value",
			Args: []string{"gravl", "--json", "foo"},
			Err:  "cannot find sample",
			Counters: map[string]int{
				"does.not.exist": 1,
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

func TestEncodeJSON(t *testing.T) {
	tests := []*internal.Harness{
		{
			Name:     "json encode",
			Args:     []string{"gravl", "--json", "encode"},
			Counters: map[string]int{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, nil, func(_ *testing.T, _ string) *cli.Command {
				return &cli.Command{
					Name: "encode",
					Action: func(c *cli.Context) error {
						return gravl.Runtime(c).Encoder.Encode(map[string]string{"hello": "world"})
					},
				}
			})
		})
	}
}

func TestHarnessHelpers(t *testing.T) {
	a := assert.New(t)
	encode := func(c *cli.Context) error { return gravl.Runtime(c).Encoder.Encode("x") }
	tests := []struct {
		tt      *internal.Harness
		handler http.Handler
	}{
		{
			// without --json the encoder discards
			tt: &internal.Harness{Name: "quiet encoder", Args: []string{"gravl", "enc"}, Action: encode},
		},
		{
			tt: &internal.Harness{
				Name: "failing encoder", Args: []string{"gravl", "enc"}, Action: encode,
				Before: internal.FailEncoding, Err: internal.ErrEncode.Error(),
			},
		},
		{
			tt: &internal.Harness{
				Name: "any error", Args: []string{"gravl", "enc"}, AnyErr: true,
				Action: func(*cli.Context) error { return errors.New("whatever it says") },
			},
		},
		{
			tt: &internal.Harness{
				Name: "stdin", Args: []string{"gravl", "enc"}, Stdin: "hello\n",
				Action: func(c *cli.Context) error {
					b, err := io.ReadAll(c.App.Reader)
					a.Equal("hello\n", string(b))
					return err
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.tt.Name, func(t *testing.T) {
			internal.Run(t, tc.tt, tc.handler, func(_ *testing.T, _ string) *cli.Command {
				return &cli.Command{Name: "enc", Action: tc.tt.Action}
			})
		})
	}
}

func TestFailingAPI(t *testing.T) {
	a := assert.New(t)
	svr := httptest.NewServer(internal.FailingAPI(t))
	defer svr.Close()
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, svr.URL, http.NoBody)
	a.NoError(err)
	res, err := http.DefaultClient.Do(req)
	a.NoError(err)
	defer res.Body.Close()
	a.Equal(http.StatusInternalServerError, res.StatusCode)
	a.Equal("application/json", res.Header.Get("Content-Type"))
	var body map[string]any
	a.NoError(json.NewDecoder(res.Body).Decode(&body))
	a.Equal("provider is down", body["message"])
}
