package route

import "encoding/json"

// Coordinates of a Route
type Coordinates [][]float64

// Routeable returns a route
type Routeable interface {
	// Route represents a series of one or more points
	Route() Coordinates
}

// Source of the Route
type Source string

// Origin of the Route
type Origin string

const (
	// Activity is a route originated from a gps track
	Activity Origin = "Activity"
	// Planned is a route originated from looking at maps
	Planned Origin = "Planned"
	// Unknown origin
	Unknown Origin = "Unknown"
)

// Route represents a series of one or more points
type Route struct {
	ID          string
	Name        string
	Source      Source
	Origin      Origin
	Description string
	Coordinates Coordinates
}

// Route returns the coordinates of the route
func (r *Route) Route() Coordinates {
	return r.Coordinates
}

// MarshalJSON produces GeoJSON
func (r *Route) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type       string                 `json:"type"`
		Geometry   map[string]interface{} `json:"geometry,omitempty"`
		Properties map[string]interface{} `json:"properties,omitempty"`
	}{
		Type: "Feature",
		Geometry: map[string]interface{}{
			"type":        "MultiPoint",
			"coordinates": r.Coordinates,
		},
		Properties: map[string]interface{}{
			"name":        r.Name,
			"source":      r.Source,
			"origin":      r.Origin,
			"description": r.Description,
		},
	})
}
