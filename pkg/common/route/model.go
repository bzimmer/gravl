package route

//go:generate stringer -type=Origin

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
type Origin int

const (
	// Activity is a route originated from a gps track
	Activity Origin = iota
	// Planned is a route originated from creating a route with a route builder
	Planned
	// Unknown origin
	Unknown
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
		ID         string                 `json:"id"`
		Type       string                 `json:"type"`
		Geometry   map[string]interface{} `json:"geometry"`
		Properties map[string]interface{} `json:"properties"`
	}{
		ID:   r.ID,
		Type: "Feature",
		Geometry: map[string]interface{}{
			"type":        "MultiPoint",
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
