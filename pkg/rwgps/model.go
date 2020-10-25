package rwgps

import "time"

// Paginator .
type Paginator struct {
	offset int
	limit  int
}

// UserResponse .
type UserResponse struct {
	User *User `json:"user"`
}

// User .
type User struct {
	AuthToken string `json:"auth_token"`
	ID        int64  `json:"id"`
	Name      string `json:"name"`
}

// TrackPoint .
type TrackPoint struct {
	Distance  float64 `json:"d"`
	Elevation float64 `json:"e"`
	Time      float64 `json:"t"`
	Longitude float64 `json:"x"`
	Latitude  float64 `json:"y"`
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
