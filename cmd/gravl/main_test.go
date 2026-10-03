package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"

	api "github.com/bzimmer/activity/strava"

	"github.com/bzimmer/gravl"
)

type fakeFaultError struct {
	code int
}

func (f *fakeFaultError) Error() string       { return "fake fault" }
func (f *fakeFaultError) HTTPStatusCode() int { return f.code }

func TestExitCode(t *testing.T) {
	for _, tt := range []struct {
		name string
		err  error
		code int
	}{
		{name: "nil error", err: nil, code: 1},
		{name: "plain error", err: errors.New("boom"), code: 1},
		{name: "non-rate-limit fault", err: &fakeFaultError{code: http.StatusInternalServerError}, code: 1},
		{name: "rate limited fault", err: &fakeFaultError{code: http.StatusTooManyRequests}, code: exitTempFail},
		{
			name: "wrapped rate limited fault",
			err:  fmt.Errorf("wrap: %w", &fakeFaultError{code: http.StatusTooManyRequests}),
			code: exitTempFail,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.code, exitCode(tt.err))
		})
	}
}

// TestExitCodeStravaRateLimit proves the full chain against a real *strava.Fault,
// not the hand-rolled fakeFaultError above: a genuine 429 HTTP response, decoded
// by github.com/bzimmer/activity's client into strava.Fault, fed into exitCode().
//
// This stops short of a real subprocess run of the gravl binary because
// strava.Before (activity/strava/strava.go) has no base-URL override -- it always
// refreshes against the live Strava OAuth endpoint -- so "gravl strava activity"
// can't be pointed at a fake server without a production code change.
func TestExitCodeStravaRateLimit(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/athlete", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"message": "Rate Limit Exceeded"}`))
	})
	svr := httptest.NewServer(mux)
	defer svr.Close()

	client, err := api.NewClient(
		api.WithBaseURL(svr.URL),
		api.WithTokenCredentials("token", "token", time.Now().Add(time.Hour)))
	assert.NoError(t, err)

	_, err = client.Athlete.Athlete(context.Background())
	assert.Error(t, err)

	assert.Equal(t, exitTempFail, exitCode(err))
}

// run runs the real application with args, returning its output and error.
func run(t *testing.T, args ...string) (string, error) {
	t.Helper()
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	app := newApp(cancel)
	var stdout, stderr bytes.Buffer
	app.Writer, app.ErrWriter = &stdout, &stderr
	app.ExitErrHandler = func(*cli.Context, error) {} // never os.Exit in a test
	err := app.RunContext(ctx, append([]string{"gravl"}, args...))
	return stdout.String(), err
}

func TestApp(t *testing.T) {
	a := assert.New(t)

	out, err := run(t, "-j", "version")
	a.NoError(err)
	a.Contains(out, `"version"`)

	_, err = run(t, "--verbosity", "loud", "version")
	a.ErrorContains(err, "loud")

	_, err = run(t, "-m", "--verbosity", "debug", "version")
	a.NoError(err)

	out, err = run(t, "-j", "qp", "providers")
	a.NoError(err)
	var providers struct {
		Exporters []string `json:"exporters"`
		Uploaders []string `json:"uploaders"`
	}
	a.NoError(json.Unmarshal([]byte(out), &providers))
	a.ElementsMatch([]string{"strava", "zwift", "hammerhead"}, providers.Exporters)
	a.ElementsMatch([]string{"strava", "cyclinganalytics"}, providers.Uploaders)

	// every command is registered
	names := make([]string, 0)
	for _, c := range commands() {
		names = append(names, c.Name)
	}
	a.Subset(names, []string{"cyclinganalytics", "hammerhead", "qp", "rwgps", "strava", "version", "zwift"})
}

func TestProviderFactories(t *testing.T) {
	a := assert.New(t)
	app := &cli.App{
		Name:     "gravl",
		Metadata: map[string]any{},
		Flags:    flags(),
		Before:   gravl.Befores(initRuntime, initQP),
		Action: func(c *cli.Context) error {
			rt := gravl.Runtime(c)
			// hammerhead's Before refreshes a token over the network, so it is left out
			for _, name := range []string{"strava", "zwift"} {
				exp, err := rt.Exporters[name](c)
				a.NoError(err, name)
				a.NotNil(exp, name)
			}
			for _, name := range []string{"strava", "cyclinganalytics"} {
				upd, err := rt.Uploaders[name](c)
				a.NoError(err, name)
				a.NotNil(upd, name)
			}
			return nil
		},
	}
	app.Writer, app.ErrWriter = io.Discard, io.Discard
	a.NoError(app.RunContext(t.Context(), []string{"gravl"}))
}

func TestEncoder(t *testing.T) {
	a := assert.New(t)
	var buf bytes.Buffer
	enc := encoder{pool: &sync.Pool{New: func() any { return json.NewEncoder(&buf) }}}
	a.NoError(enc.Encode(map[string]int{"a": 1}))
	a.NoError(enc.Encode(map[string]int{"b": 2}))
	a.Equal("{\"a\":1}\n{\"b\":2}\n", buf.String())

	// a value json cannot encode returns the error and drops the encoder
	a.Error(enc.Encode(make(chan int)))

	a.NoError(encoder{}.Encode("discarded"))

	wrong := encoder{pool: &sync.Pool{New: func() any { return "not an encoder" }}}
	a.ErrorContains(wrong.Encode(1), "did not receive encoder")
}

func TestSignalCancels(t *testing.T) {
	a := assert.New(t)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	cancelled := make(chan struct{})
	app := &cli.App{
		Name:   "gravl",
		Before: initSignal(func() { close(cancelled) }),
		Action: func(*cli.Context) error {
			// an interrupt reaches the handler, which cancels
			p, err := os.FindProcess(os.Getpid())
			if err != nil {
				return err
			}
			if err = p.Signal(os.Interrupt); err != nil {
				return err
			}
			select {
			case <-cancelled:
				return nil
			case <-time.After(5 * time.Second):
				return errors.New("interrupt did not cancel")
			}
		},
	}
	a.NoError(app.RunContext(ctx, []string{"gravl"}))
}
