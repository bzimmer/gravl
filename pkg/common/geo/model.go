package geo

import (
	"encoding/json"

	geom "github.com/twpayne/go-geom"
	gpx "github.com/twpayne/go-gpx"
)

type GPX interface {
	GPX() (*gpx.GPX, error)
}

// GeographicName information about the official name for places, features, and areas
type GeographicName struct {
	ID          string
	Name        string
	Source      string
	Description string
	Coordinates geom.Coord
	Class       string
	Locale      string
}

// MarshalJSON produces GeoJSON
func (g *GeographicName) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID         string                 `json:"id"`
		Type       string                 `json:"type"`
		Geometry   map[string]interface{} `json:"geometry"`
		Properties map[string]interface{} `json:"properties"`
	}{
		ID:   g.ID,
		Type: "Feature",
		Geometry: map[string]interface{}{
			"type":        "MultiPoint",
			"coordinates": g.Coordinates,
		},
		Properties: map[string]interface{}{
			"name":        g.Name,
			"source":      g.Source,
			"class":       g.Class,
			"locale":      g.Locale,
			"description": g.Description,
		},
	})
}
