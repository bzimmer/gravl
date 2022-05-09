package internal

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/armon/go-metrics"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"

	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/eval/antonmedv"
)

type Harness struct {
	Name, Err string
	Args      []string
	Counters  map[string]int
	Before    cli.BeforeFunc
	After     cli.AfterFunc
	Action    cli.ActionFunc
}

func runtime(app *cli.App) *pkg.Rt {
	return app.Metadata[pkg.RuntimeKey].(*pkg.Rt)
}

func initRuntime(c *cli.Context) error {
	var enc pkg.Encoder
	cfg := metrics.DefaultConfig("gravl")
	cfg.EnableRuntimeMetrics = false
	cfg.TimerGranularity = time.Second
	sink := metrics.NewInmemSink(time.Hour*24, time.Hour*24)
	metric, err := metrics.New(cfg, sink)
	if err != nil {
		return err
	}
	switch c.String("encoding") {
	case "json":
		enc = pkg.JSON(c.App.Writer, false)
	default:
		enc = pkg.Blackhole()
	}
	c.App.Metadata = map[string]any{
		pkg.RuntimeKey: &pkg.Rt{
			Start:     time.Now(),
			Metrics:   metric,
			Sink:      sink,
			Encoder:   enc,
			Fs:        afero.NewMemMapFs(),
			Filterer:  antonmedv.Filterer,
			Evaluator: antonmedv.Evaluator,
			Exporters: make(map[string]pkg.ExporterFunc),
			Uploaders: make(map[string]pkg.UploaderFunc),
			Endpoints: make(map[string]oauth2.Endpoint),
		},
	}
	log.Info().Msg("initiated Runtime")
	return nil
}

func counters(t *testing.T, expected map[string]int) cli.AfterFunc {
	a := assert.New(t)
	return func(c *cli.Context) error {
		data := pkg.Runtime(c).Sink.Data()
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
		Before:   pkg.Befores(initRuntime, tt.Before),
		After:    pkg.Afters(tt.After, walkfs, pkg.Stats, counters(t, tt.Counters)),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "encoding",
				Aliases: []string{"e"},
				Value:   "",
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
