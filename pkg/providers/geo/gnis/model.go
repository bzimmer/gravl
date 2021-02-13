package gnis

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"

	"github.com/bzimmer/gravl/pkg/providers/geo"
)

var _ geo.GeoJSON = &GeographicName{}

// GeographicName information about the official name for places, features, and areas
type GeographicName struct {
	ID     string      `json:"id,omitempty"`
	Name   string      `json:"name,omitempty"`
	Source string      `json:"source,omitempty"`
	Class  string      `json:"class,omitempty"`
	Locale string      `json:"locale,omitempty"`
	Point  *geom.Point `json:"point,omitempty"`
}

func (g *GeographicName) MarshalJSON() ([]byte, error) {
	f, err := g.GeoJSON()
	if err != nil {
		return nil, err
	}
	return f.MarshalJSON()
}

func (g *GeographicName) GeoJSON() (*geojson.Feature, error) {
	return &geojson.Feature{
		ID:       g.ID,
		Geometry: g.Point,
		Properties: map[string]interface{}{
			"name":   g.Name,
			"source": g.Source,
			"class":  g.Class,
			"locale": g.Locale,
		},
	}, nil
}
