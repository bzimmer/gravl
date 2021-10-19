package pkg

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"time"

	"github.com/armon/go-metrics"
	"github.com/bzimmer/activity"
	"github.com/bzimmer/activity/cyclinganalytics"
	"github.com/bzimmer/activity/rwgps"
	"github.com/bzimmer/activity/strava"
	"github.com/bzimmer/activity/zwift"
	"github.com/bzimmer/gravl/pkg/eval"
	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog/log"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
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

type Rt struct {
	// Metadata
	Start time.Time

	// Activity clients
	Zwift            *zwift.Client
	Strava           *strava.Client
	RideWithGPS      *rwgps.Client
	CyclingAnalytics *cyclinganalytics.Client

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
	Mapper    func(string) (eval.Mapper, error)
	Filterer  func(string) (eval.Filterer, error)
	Evaluator func(string) (eval.Evaluator, error)
}

func Runtime(c *cli.Context) *Rt {
	return c.App.Metadata[RuntimeKey].(*Rt)
}

var ErrUnknownEncoder = errors.New("unknown encoder")

type Encoder interface {
	Encode(v interface{}) error
}

type blackhole struct{}

func (b *blackhole) Encode(v interface{}) error {
	return nil
}

type spewEncoder struct {
	cfg    *spew.ConfigState
	writer io.Writer
}

func (s *spewEncoder) Encode(v interface{}) error {
	s.cfg.Fdump(s.writer, v)
	return nil
}

type gpxEncoder struct {
	enc Encoder
}

func (g *gpxEncoder) Encode(v interface{}) error {
	q, ok := v.(activity.GPXer)
	if !ok {
		return errors.New("encoding GPX not supported")
	}
	v, err := q.GPX()
	if err != nil {
		return err
	}
	return g.enc.Encode(v)
}

type geoJSONEncoder struct {
	enc Encoder
}

func (g *geoJSONEncoder) Encode(v interface{}) error {
	q, ok := v.(activity.GeoJSONer)
	if !ok {
		return errors.New("encoding GeoJSON not supported")
	}
	v, err := q.GeoJSON()
	if err != nil {
		return err
	}
	return g.enc.Encode(v)
}

func (g *geoJSONEncoder) Name() string {
	return "geojson"
}

type jsonEncoder struct {
	enc *json.Encoder
}

func (j *jsonEncoder) Encode(v interface{}) error {
	return j.enc.Encode(v)
}

type xmlEncoder struct {
	enc *xml.Encoder
}

func (x *xmlEncoder) Encode(v interface{}) error {
	return x.enc.Encode(v)
}

func JSON(writer io.Writer, compact bool) Encoder {
	enc := json.NewEncoder(writer)
	if !compact {
		enc.SetIndent("", " ")
	}
	enc.SetEscapeHTML(false)
	return &jsonEncoder{enc: enc}
}

func XML(writer io.Writer, compact bool) Encoder {
	enc := xml.NewEncoder(writer)
	if !compact {
		enc.Indent("", " ")
	}
	return &xmlEncoder{enc: enc}
}

func GeoJSON(writer io.Writer, compact bool) Encoder {
	return &geoJSONEncoder{enc: JSON(writer, compact)}
}

func GPX(writer io.Writer, compact bool) Encoder {
	return &gpxEncoder{enc: XML(writer, compact)}
}

func Spew(writer io.Writer) Encoder {
	cfg := spew.NewDefaultConfig()
	cfg.SortKeys = true
	return &spewEncoder{cfg: cfg, writer: writer}
}

func Blackhole() Encoder {
	return &blackhole{}
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
