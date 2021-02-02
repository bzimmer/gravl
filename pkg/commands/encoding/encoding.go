package encoding

//go:generate stringer -type=Encoding -linecomment -output=encoding_string.go

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"

	"github.com/bzimmer/gravl/pkg/providers/geo"
	"github.com/davecgh/go-spew/spew"
)

type EncodeFunc func(v interface{}) error

var ErrUnknownEncoder = errors.New("unknown encoder")
var ErrExistingEncoder = errors.New("encoder exists")

type Encoder interface {
	Name() string
	Encode(v interface{}) error
}

var Encode EncodeFunc = func(v interface{}) error {
	return nil
}

type Encoders struct {
	encoders map[string]Encoder
}

// Use the encoder for the encoding if no prior encoder exists
func (n *Encoders) Use(encoder Encoder) error {
	_, ok := n.encoders[encoder.Name()]
	if ok {
		return ErrExistingEncoder
	}
	n.encoders[encoder.Name()] = encoder
	return nil
}

// MustUse the encoder for the encoding
func (n *Encoders) MustUse(encoder Encoder) {
	_ = n.Use(encoder)
}

// For name return an encoder
func (n *Encoders) For(encoder string) (Encoder, error) {
	enc, ok := n.encoders[encoder]
	if !ok {
		return nil, ErrUnknownEncoder
	}
	return enc, nil
}

func NewEncoders() *Encoders {
	return &Encoders{encoders: make(map[string]Encoder)}
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
	if q, ok := v.(geo.GPX); ok {
		p, err := q.GPX()
		if err != nil {
			return err
		}
		v = p
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
	if q, ok := v.(geo.GeoJSON); ok {
		p, err := q.GeoJSON()
		if err != nil {
			return err
		}
		v = p
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
