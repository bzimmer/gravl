package encoding

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"

	"github.com/davecgh/go-spew/spew"

	"github.com/bzimmer/gravl/pkg/providers/activity"
	"github.com/bzimmer/gravl/pkg/providers/geo"
)

var ErrUnknownEncoder = errors.New("unknown encoder")
var ErrExistingEncoder = errors.New("encoder exists")

type Encoder interface {
	Name() string
	Encode(v interface{}) error
}

var DefaultEncoder = "json"
var encoders = make(map[string]Encoder)

var Encode = func(v interface{}) error {
	enc, err := For(DefaultEncoder)
	if err != nil {
		return err
	}
	return enc.Encode(v)
}

// Add the encoder for the encoding if no prior encoder exists
//
// This function is not thread-safe
func Add(encoder Encoder) {
	encoders[encoder.Name()] = encoder
}

// For name return an encoder
func For(encoder string) (Encoder, error) {
	enc, ok := encoders[encoder]
	if !ok {
		return nil, ErrUnknownEncoder
	}
	return enc, nil
}

type namedEncoder struct {
	enc Encoder
}

func (x *namedEncoder) Name() string {
	return "named"
}

func (x *namedEncoder) Encode(v interface{}) error {
	if n, ok := v.(activity.Named); ok {
		return x.enc.Encode(n.Handle())
	}
	return errors.New("not Named")
}

type spewEncoder struct {
	cfg    *spew.ConfigState
	writer io.Writer
}

func (s *spewEncoder) Encode(v interface{}) error {
	s.cfg.Fdump(s.writer, v)
	return nil
}

func (s *spewEncoder) Name() string {
	return "spew"
}

type gpxEncoder struct {
	enc Encoder
}

func (g *gpxEncoder) Encode(v interface{}) error {
	q, ok := v.(geo.GPX)
	if !ok {
		return errors.New("encoding GPX not supported")
	}
	v, err := q.GPX()
	if err != nil {
		return err
	}
	return g.enc.Encode(v)
}

func (g *gpxEncoder) Name() string {
	return "gpx"
}

type geoJSONEncoder struct {
	enc Encoder
}

func (g *geoJSONEncoder) Encode(v interface{}) error {
	q, ok := v.(geo.GeoJSON)
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

func (j *jsonEncoder) Name() string {
	return "json"
}

type xmlEncoder struct {
	enc *xml.Encoder
}

func (x *xmlEncoder) Encode(v interface{}) error {
	return x.enc.Encode(v)
}

func (x *xmlEncoder) Name() string {
	return "xml"
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

func Named(writer io.Writer, compact bool) Encoder {
	return &namedEncoder{enc: JSON(writer, compact)}
}
