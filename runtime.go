package gravl

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/bzimmer/activity"
	"github.com/bzimmer/activity/cyclinganalytics"
	"github.com/bzimmer/activity/rwgps"
	"github.com/bzimmer/activity/strava"
	"github.com/bzimmer/activity/zwift"
	"github.com/hashicorp/go-metrics"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"

	"github.com/bzimmer/gravl/eval"
)

const RuntimeKey = "github.com/bzimmer/gravl#RuntimeKey"

// Token produces a random token of length `n`
func Token(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

type ExporterFunc func(c *cli.Context) (activity.Exporter, error)
type UploaderFunc func(c *cli.Context) (activity.Uploader, error)

// Rt holds the gravl runtime
type Rt struct {
	// Metadata
	Start time.Time

	// Activity clients
	Zwift            *zwift.Client
	Strava           *strava.Client
	RideWithGPS      *rwgps.Client
	CyclingAnalytics *cyclinganalytics.Client

	// Endpoints
	Endpoints map[string]oauth2.Endpoint

	// Export / Upload
	Exporters map[string]ExporterFunc
	Uploaders map[string]UploaderFunc

	// IO
	Fs      afero.Fs
	Encoder Encoder

	// Metrics
	Metrics *metrics.Metrics
	Sink    *metrics.InmemSink

	// Evaluation
	Filterer  func(string) (eval.Filterer, error)
	Evaluator func(string) (eval.Evaluator, error)
}

func Runtime(c *cli.Context) *Rt {
	return c.App.Metadata[RuntimeKey].(*Rt) //nolint:errcheck // cannot happen
}

type Encoder interface {
	Encode(v any) error
}

// Afters combines multiple `cli.AfterFunc`s into a single `cli.AfterFunc`
func Afters(afs ...cli.AfterFunc) cli.AfterFunc {
	return func(c *cli.Context) error {
		for _, fn := range afs {
			if fn == nil {
				continue
			}
			if err := fn(c); err != nil {
				return err
			}
		}
		return nil
	}
}

// Befores combines multiple `cli.BeforeFunc`s into a single `cli.BeforeFunc`
func Befores(bfs ...cli.BeforeFunc) cli.BeforeFunc {
	return func(c *cli.Context) error {
		for _, fn := range bfs {
			if fn == nil {
				continue
			}
			if err := fn(c); err != nil {
				return err
			}
		}
		return nil
	}
}

// Stats logs and encodes (if requested) the stats
func Stats(c *cli.Context) error {
	data := Runtime(c).Sink.Data()
	for i := range data {
		for key, val := range data[i].Counters {
			log.Info().
				Int("count", val.Count).
				Str("metric", key).
				Msg("counters")
		}
		for key, val := range data[i].Samples {
			as := val.AggregateSample
			log.Info().
				Int("count", val.Count).
				Str("metric", key).
				Float64("min", as.Min).
				Float64("max", as.Max).
				Float64("mean", as.Mean()).
				Float64("stddev", as.Stddev()).
				Msg("samples")
		}
	}
	return nil
}
