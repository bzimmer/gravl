package rwgps

//go:generate stringer -type=Origin

import (
	"fmt"
	"time"

	"github.com/bzimmer/gravl/pkg/common/geo"
)

// Origin of the trip
type Origin int

const (
	// TripOrigin is a ride which was recorded by GPS
	TripOrigin Origin = iota
	// RouteOrigin is a ride which was planned on the RWGPS route builder
	RouteOrigin
)

// UserResponse .
type UserResponse struct {
	User *User `json:"user"`
}

// User .
type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	AuthToken string `json:"auth_token"`
}

// TrackPoint .
type TrackPoint struct {
	Longitude float64 `json:"x"`
	Latitude  float64 `json:"y"`
	Elevation float64 `json:"e"` // elevation in meters
	Distance  float64 `json:"d"` // distance in meters
	Time      float64 `json:"t"` // seconds since epoch, unix timestamp
}

// Trip .
type Trip struct {
	CreatedAt     time.Time     `json:"created_at"`
	DepartedAt    time.Time     `json:"departed_at"`
	Description   string        `json:"description"`
	Distance      float64       `json:"distance"`
	ElevationGain float64       `json:"elevation_gain"`
	ElevationLoss float64       `json:"elevation_loss"`
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	Origin        Origin        `json:"-"`
	TrackID       string        `json:"track_id"`
	TrackPoints   []*TrackPoint `json:"track_points"`
	UpdatedAt     time.Time     `json:"updated_at"`
	UserID        int           `json:"user_id"`
	Visibility    int           `json:"visibility"`
}

// TripResponse .
type TripResponse struct {
	Type  string `json:"type"`
	Trip  *Trip  `json:"trip"`
	Route *Trip  `json:"route"`
}

func (t *Trip) Track() (*geo.Track, error) {
	coords := make([][]float64, len(t.TrackPoints))
	for i, tp := range t.TrackPoints {
		coords[i] = []float64{tp.Longitude, tp.Latitude, tp.Elevation}
	}

	var q geo.Origin
	switch t.Origin {
	case TripOrigin:
		q = geo.Activity
	case RouteOrigin:
		q = geo.Planned
	}

	return &geo.Track{
		ID:          fmt.Sprintf("%d", t.ID),
		Name:        t.Name,
		Source:      baseURL,
		Origin:      q,
		Description: t.Description,
		Coordinates: coords,
	}, nil
}
