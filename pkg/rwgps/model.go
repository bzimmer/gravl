package rwgps

import "time"

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
