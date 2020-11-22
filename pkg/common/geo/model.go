package geo

//go:generate stringer -type=Origin

import "encoding/json"

// Coordinates of a Route
type Coordinates [][]float64

// Source of the Route
type Source string

// Origin of the Route
type Origin int

const (
	// Activity is a route originated from a gps track
	Activity Origin = iota
	// Planned is a route originated from creating a route with a route builder
	Planned
	// Unknown origin
	Unknown
)

// Track represents a series of one or more points
type Track struct {
	ID          string
	Name        string
	Source      Source
	Origin      Origin
	Description string
	Coordinates Coordinates
}

// Trackable instances can return a Track
type Trackable interface {
	// Track returns an instance of a Track
	Track() (*Track, error)
}

// MarshalJSON produces GeoJSON
func (r *Track) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID         string                 `json:"id"`
		Type       string                 `json:"type"`
		Geometry   map[string]interface{} `json:"geometry"`
		Properties map[string]interface{} `json:"properties"`
	}{
		ID:   r.ID,
		Type: "Feature",
		Geometry: map[string]interface{}{
			"type":        "LineString",
			"coordinates": r.Coordinates,
		},
		Properties: map[string]interface{}{
			"name":        r.Name,
			"source":      r.Source,
			"origin":      r.Origin.String(),
			"description": r.Description,
		},
	})
}

// GeographicName information about the official name for places, features, and areas
type GeographicName struct {
	Track
	Class  string
	Locale string
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
