package gravl

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"os"

	"github.com/bzimmer/gravl/pkg/common/geo"
)

type Encoding int

const (
	EncodingNative Encoding = iota // native
	EncodingXML                    // xml
	EncodingJSON                   // json
)

type xcoder struct {
	enc  Encoding
	xml  *xml.Encoder
	json *json.Encoder
}

func (x *xcoder) Encode(v interface{}) error {
	switch x.enc {
	case EncodingXML:
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
	}
	return nil
}

func newEncoder(writer io.Writer, encoding string, compact bool) *xcoder {
	if writer == nil {
		writer = os.Stdout
	}
	xe := xml.NewEncoder(writer)
	xe.Indent("", " ")
	je := json.NewEncoder(writer)
	if !compact {
		je.SetIndent("", " ")
	}
	je.SetEscapeHTML(false)

	var enc Encoding
	switch encoding {
	case "native":
		enc = EncodingNative
	case "json", "geojson":
		enc = EncodingJSON
	case "xml", "gpx":
		enc = EncodingXML
	}
	return &xcoder{
		enc:  enc,
		xml:  xe,
		json: je,
	}
}
