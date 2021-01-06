package encoding

//go:generate stringer -type=Encoding -linecomment -output=encoding_string.go

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/bzimmer/gravl/pkg/geo"
)

type Encoding int

type Encoder func(v interface{}) error

const (
	EncodingNative  Encoding = iota // native
	EncodingXML                     // xml
	EncodingJSON                    // json
	EncodingGeoJSON                 // geojson
	EncodingGPX                     // gpx
	EncodingSpew                    // spew
)

func Encode(v interface{}) error {
	x, err := NewEncoder(os.Stdout, EncodingNative.String(), false)
	if err != nil {
		return err
	}
	return x.Encode(v)
}

type Xcoder struct {
	enc             Encoding
	xml, json, spew Encoder
}

func (x *Xcoder) Encode(v interface{}) error {
	switch x.enc {
	case EncodingXML:
		return x.xml(v)
	case EncodingGPX:
		if q, ok := v.(geo.GPX); ok {
			p, err := q.GPX()
			if err != nil {
				return err
			}
			v = p
		}
		return x.xml(v)
	case EncodingNative, EncodingJSON:
		return x.json(v)
	case EncodingGeoJSON:
		if q, ok := v.(geo.GeoJSON); ok {
			p, err := q.GeoJSON()
			if err != nil {
				return err
			}
			v = p
		}
		return x.json(v)
	case EncodingSpew:
		return x.spew(v)
	}
	return nil
}

func NewEncoder(writer io.Writer, encoding string, compact bool) (*Xcoder, error) {
	if writer == nil {
		writer = os.Stdout
	}
	xe := xml.NewEncoder(writer)
	if !compact {
		xe.Indent("", " ")
	}
	je := json.NewEncoder(writer)
	if !compact {
		je.SetIndent("", " ")
	}
	je.SetEscapeHTML(false)

	cfg := spew.NewDefaultConfig()
	cfg.SortKeys = true
	var sw = func(v interface{}) (err error) {
		cfg.Fdump(writer, v)
		return
	}

	var enc Encoding
	switch encoding {
	case "native":
		enc = EncodingNative
	case "json":
		enc = EncodingJSON
	case "geojson":
		enc = EncodingGeoJSON
	case "xml":
		enc = EncodingXML
	case "gpx":
		enc = EncodingGPX
	case "spew", "dump":
		enc = EncodingSpew
	default:
		return nil, fmt.Errorf("unknown encoder: '%s'", encoding)
	}

	return &Xcoder{enc: enc, xml: xe.Encode, json: je.Encode, spew: sw}, nil
}
