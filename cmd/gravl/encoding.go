package gravl

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/bzimmer/gravl/pkg/common/geo"
)

type Encoding int

const (
	EncodingNative  Encoding = iota // native
	EncodingXML                     // xml
	EncodingJSON                    // json
	EncodingGeoJSON                 // geojson
	EncodingGPX                     // gpx
)

type xcoder struct {
	enc  Encoding
	xml  *xml.Encoder
	json *json.Encoder
}

func (x *xcoder) Encode(v interface{}) error {
	switch x.enc {
	case EncodingXML:
		return x.xml.Encode(v)
	case EncodingGPX:
		if q, ok := v.(geo.GPX); ok {
			p, err := q.GPX()
			if err != nil {
				return err
			}
			v = p
		}
		return x.xml.Encode(v)
	case EncodingNative, EncodingJSON:
		return x.json.Encode(v)
	case EncodingGeoJSON:
		if q, ok := v.(geo.GeoJSON); ok {
			p, err := q.GeoJSON()
			if err != nil {
				return err
			}
			v = p
		}
		return x.json.Encode(v)
	}
	return nil
}

func newEncoder(writer io.Writer, encoding string, compact bool) (*xcoder, error) {
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
	default:
		return nil, fmt.Errorf("unknown encoder: %s", encoding)
	}
	return &xcoder{enc: enc, xml: xe, json: je}, nil
}
