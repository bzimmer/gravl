package internal

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/armon/go-metrics"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/eval/antonmedv"
)

type Harness struct {
	Name, Err string
	Args      []string
	Counters  map[string]int
	Before    cli.BeforeFunc
	After     cli.AfterFunc
	Action    cli.ActionFunc
}

func runtime(app *cli.App) *gravl.Rt {
	return app.Metadata[gravl.RuntimeKey].(*gravl.Rt)
}

type encoder struct {
	pool *sync.Pool
}

func (enc encoder) Encode(v any) error {
	if enc.pool == nil {
		return nil
	}
	x, ok := enc.pool.Get().(*json.Encoder)
	if !ok {
		return errors.New("did not receive encoder from pool")
	}
	defer enc.pool.Put(x)
	return x.Encode(v)
}

func initRuntime(c *cli.Context) error {
	cfg := metrics.DefaultConfig("gravl")
	cfg.EnableRuntimeMetrics = false
	cfg.TimerGranularity = time.Second
	sink := metrics.NewInmemSink(time.Hour*24, time.Hour*24)
	metric, err := metrics.New(cfg, sink)
	if err != nil {
		return err
	}
	var pool *sync.Pool
	if c.Bool("json") {
		pool = &sync.Pool{New: func() any { return json.NewEncoder(c.App.Writer) }}
	}
	c.App.Metadata = map[string]any{
		gravl.RuntimeKey: &gravl.Rt{
			Start:     time.Now(),
			Metrics:   metric,
			Sink:      sink,
			Encoder:   &encoder{pool: pool},
			Fs:        afero.NewMemMapFs(),
			Filterer:  antonmedv.Filterer,
			Evaluator: antonmedv.Evaluator,
			Exporters: make(map[string]gravl.ExporterFunc),
			Uploaders: make(map[string]gravl.UploaderFunc),
			Endpoints: make(map[string]oauth2.Endpoint),
		},
	}
	log.Info().Msg("initiated Runtime")
	return nil
}

func counters(t *testing.T, expected map[string]int) cli.AfterFunc {
	a := assert.New(t)
	return func(c *cli.Context) error {
		data := gravl.Runtime(c).Sink.Data()
		for key, value := range expected {
			var found bool
			for i := range data {
				if counter, ok := data[i].Counters[key]; ok {
					found = true
					a.Equalf(value, counter.Count, key)
					break
				}
			}
			if !found {
				return fmt.Errorf("cannot find sample value for {%s}", key)
			}
		}
		return nil
	}
}

func walkfs(c *cli.Context) error {
	return afero.Walk(runtime(c.App).Fs, "/", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Fprintf(c.App.ErrWriter, "%s\n", path)
		return nil
	})
}

func TestMain(m *testing.M) {
	// hijack the `go test` verbose flag to manage logging
	verbose := flag.CommandLine.Lookup("test.v")
	if verbose.Value.String() != "" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	os.Exit(m.Run())
}

func Run(t *testing.T, tt *Harness, handler http.Handler, cmd func(*testing.T, string) *cli.Command) {
	RunContext(context.Background(), t, tt, handler, cmd)
}

func RunContext(ctx context.Context, t *testing.T, tt *Harness, handler http.Handler, cmd func(*testing.T, string) *cli.Command) {
	a := assert.New(t)

	svr := httptest.NewServer(handler)
	defer svr.Close()

	app := NewTestApp(t, tt, cmd(t, svr.URL))
	err := app.RunContext(ctx, tt.Args)
	switch tt.Err == "" {
	case true:
		a.NoError(err)
	case false:
		a.Error(err)
		a.Contains(err.Error(), tt.Err)
	}
}

func NewTestApp(t *testing.T, tt *Harness, cmd *cli.Command) *cli.App {
	return &cli.App{
		Name:     tt.Name,
		HelpName: tt.Name,
		Before:   gravl.Befores(initRuntime, tt.Before),
		After:    gravl.Afters(tt.After, walkfs, gravl.Stats, counters(t, tt.Counters)),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "json",
				Aliases: []string{"j"},
				Value:   false,
			},
			&cli.BoolFlag{
				Name:  "http-tracing",
				Value: false,
				Usage: "Log all http calls (warning: no effort is made to mask log ids, keys, and other sensitive information)",
			},
			&cli.DurationFlag{
				Name:    "timeout",
				Aliases: []string{"t"},
				Value:   time.Second * 10,
			}},
		Commands: []*cli.Command{cmd},
	}
}
