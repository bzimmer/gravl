package internal

import (
	"context"
	"flag"
	"fmt"
	"io"
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
}

func runtime(app *cli.App) *pkg.Rt {
	return app.Metadata[pkg.RuntimeKey].(*pkg.Rt)
}

func before(c *cli.Context) error {
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
	c.App.Metadata = map[string]interface{}{
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
	return nil
}

func CopyFile(w io.Writer, filename string) error {
	fp, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = io.Copy(w, fp)
	return err
}

func findCounter(app *cli.App, name string) (metrics.SampledValue, error) {
	sink := runtime(app).Sink
	for i := range sink.Data() {
		im := sink.Data()[i]
		if sample, ok := im.Counters[name]; ok {
			return sample, nil
		}
	}
	return metrics.SampledValue{}, fmt.Errorf("cannot find sample value for {%s}", name)
}

func TestMain(m *testing.M) {
	// hijack the `go test` verbose flag to manage logging
	verbose := flag.CommandLine.Lookup("test.v")
	if verbose.Value.String() != "" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	os.Exit(m.Run())
}

func Run(t *testing.T, tt *Harness, mux *http.ServeMux, cmd func(*testing.T, string) *cli.Command) {
	a := assert.New(t)

	svr := httptest.NewServer(mux)
	defer svr.Close()

	app := NewTestApp(t, tt.Name, cmd(t, svr.URL))
	if tt.Before != nil {
		app.Before = pkg.Befores(app.Before, tt.Before)
	}
	if tt.After != nil {
		app.After = pkg.Afters(app.After, tt.After)
	}

	err := app.RunContext(context.Background(), tt.Args)
	switch tt.Err == "" {
	case true:
		a.NoError(err)
	case false:
		a.Error(err)
		a.Contains(err.Error(), tt.Err)
	}

	for key, value := range tt.Counters {
		counter, err := findCounter(app, key)
		a.NoError(err)
		a.Equalf(value, counter.Count, key)
	}
}

func NewTestApp(t *testing.T, name string, cmd *cli.Command) *cli.App {
	return &cli.App{
		Name:     name,
		HelpName: name,
		Before:   before,
		After: func(c *cli.Context) error {
			t.Log(name)
			if err := afero.Walk(runtime(c.App).Fs, "/", func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				fmt.Fprintf(c.App.ErrWriter, "%s\n", path)
				return nil
			}); err != nil {
				return err
			}
			return pkg.Stats(c)
		},
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
		ExitErrHandler: func(c *cli.Context, err error) {
			if err == nil {
				return
			}
			log.Error().Err(err).Msg(c.App.Name)
		},
		Commands: []*cli.Command{cmd},
	}
}
