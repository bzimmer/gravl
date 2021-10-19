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
	Before, After func(c *cli.Context) error
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
			Mapper:    antonmedv.Mapper,
			Filterer:  antonmedv.Filterer,
			Evaluator: antonmedv.Evaluator,
			Exporters: make(map[string]pkg.ExporterFunc),
			Uploaders: make(map[string]pkg.UploaderFunc),
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
		app.After = pkg.Afters(app.After, tt.After)
	}
}

func NewTestApp(t *testing.T, name string, cmd *cli.Command) *cli.App {
	return &cli.App{
		Name:     name,
		HelpName: name,
		Before:   before,
		After: func(c *cli.Context) error {
			t.Log(name)
			switch v := runtime(c.App).Fs.(type) {
			case *afero.MemMapFs:
				v.List()
			default:
			}
			return pkg.Stats(c)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "encoding",
				Aliases: []string{"e"},
				Value:   "",
				Usage:   "Output encoding (eg: json, xml, geojson, gpx, spew)",
			},
			&cli.DurationFlag{
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
	}
}
