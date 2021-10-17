package internal

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/armon/go-metrics"
	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/eval/antonmedv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

type Harness struct {
	Name, Err     string
	Args          []string
	Counters      map[string]int
	Before, After func(app *cli.App)
}

func runtime(app *cli.App) *pkg.Rt {
	return app.Metadata[pkg.RuntimeKey].(*pkg.Rt)
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
		tt.Before(app)
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

	if tt.After != nil {
		tt.After(app)
	}
}

func NewTestApp(t *testing.T, name string, cmd *cli.Command) *cli.App {
	cfg := metrics.DefaultConfig("gravl")
	cfg.EnableRuntimeMetrics = false
	cfg.TimerGranularity = time.Second
	sink := metrics.NewInmemSink(time.Hour*24, time.Hour*24)
	metric, err := metrics.New(cfg, sink)
	if err != nil {
		t.Error(err)
	}

	return &cli.App{
		Name:     name,
		HelpName: name,
		After: func(c *cli.Context) error {
			t.Log(name)
			switch v := runtime(c.App).Fs.(type) {
			case *afero.MemMapFs:
				v.List()
			default:
			}
			return pkg.Stats(c)
		},
		Flags: []cli.Flag{&cli.DurationFlag{
			Name:    "timeout",
			Aliases: []string{"t"},
			Value:   time.Second * 10,
			Usage:   "Timeout duration (eg, 1ms, 2s, 5m, 3h)",
		}},
		ExitErrHandler: func(c *cli.Context, err error) {
			if err == nil {
				return
			}
			log.Error().Err(err).Msg(c.App.Name)
		},
		Commands: []*cli.Command{cmd},
		Metadata: map[string]interface{}{
			pkg.RuntimeKey: &pkg.Rt{
				Start:     time.Now(),
				Metrics:   metric,
				Sink:      sink,
				Encoder:   pkg.Blackhole(),
				Fs:        afero.NewMemMapFs(),
				Mapper:    antonmedv.Mapper,
				Filterer:  antonmedv.Filterer,
				Evaluator: antonmedv.Evaluator,
			},
		},
	}
}
