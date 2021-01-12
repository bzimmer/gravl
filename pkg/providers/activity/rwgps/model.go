package rwgps

//go:generate stringer -type=Type -linecomment -output=model_string.go

import (
	"time"

	"github.com/martinlindhe/unit"
)

// Type of the trip
type Type int

const (
	// TypeTrip is a ride which was recorded by GPS
	TypeTrip Type = iota // trip
	// TypeRoute is a ride which was planned on the RWGPS route builder
	TypeRoute // route
)

type UserID int64

const (
	Me UserID = 0
)

// Fault .
type Fault struct {
	Message string `json:"message"`
}

func (f *Fault) Error() string {
	return f.Message
}

// UserResponse .
type UserResponse struct {
	User *User `json:"user"`
}

// User .
type User struct {
	ID        UserID `json:"id"`
	Name      string `json:"name"`
	AuthToken string `json:"auth_token"`
}

type Summary struct {
	Avg    float64 `json:"avg"`
	AvgRaw float64 `json:"_avg"`
	Max    float64 `json:"max"`
	MaxI   float64 `json:"max_i"`
	MaxRaw float64 `json:"_max"`
	Min    float64 `json:"min"`
	MinI   float64 `json:"min_i"`
	MinRaw float64 `json:"_min"`
}

type Metrics struct {
	AscentTime    int         `json:"ascentTime"`
	Cadence       *Summary    `json:"cad"`
	Calories      int         `json:"calories"`
	CreatedAt     time.Time   `json:"created_at"`
	DescentTime   int         `json:"descentTime"`
	Distance      unit.Length `json:"distance" units:"m"`
	Duration      int         `json:"duration"`
	Elevation     *Summary    `json:"ele"`
	ElevationGain unit.Length `json:"ele_gain" units:"m"`
	ElevationLoss unit.Length `json:"ele_loss" units:"m"`
	EndElevation  unit.Length `json:"endElevation" units:"m"`
	FirstTime     int         `json:"firstTime"`
	Grade         *Summary    `json:"grade"`
	HeartRate     *Summary    `json:"hr"`
	ID            int         `json:"id"`
	MovingPace    float64     `json:"movingPace"`
	MovingTime    int         `json:"movingTime"`
	NumPoints     int         `json:"numPoints"`
	Pace          float64     `json:"pace"`
	ParentID      int         `json:"parent_id"`
	ParentType    string      `json:"parent_type"`
	Speed         *Summary    `json:"speed"`
	Stationary    bool        `json:"stationary"`
	StoppedTime   int         `json:"stoppedTime"`
	UpdatedAt     *time.Time  `json:"updated_at"`
	V             int         `json:"v"`
	VAM           float64     `json:"vam"`
	Watts         *Summary    `json:"watts"`
}

// TrackPoint .
type TrackPoint struct {
	Longitude float64     `json:"x"`
	Latitude  float64     `json:"y"`
	Elevation unit.Length `json:"e" units:"m"`
	Distance  unit.Length `json:"d" units:"m"`
	Time      float64     `json:"t"` // seconds since epoch, unix timestamp
	Cadence   float64     `json:"c"`
	Grade     float64     `json:"g"`
	Speed     unit.Speed  `json:"s" units:"kph"`
}

// Trip .
type Trip struct {
	CreatedAt     time.Time     `json:"created_at"`
	DepartedAt    time.Time     `json:"departed_at"`
	Description   string        `json:"description"`
	Distance      unit.Length   `json:"distance" units:"m"`
	Duration      int           `json:"duration"`
	ElevationGain unit.Length   `json:"elevation_gain" units:"m"`
	ElevationLoss unit.Length   `json:"elevation_loss" units:"m"`
	ID            int64         `json:"id"`
	Name          string        `json:"name"`
	Type          string        `json:"type"`
	TrackID       string        `json:"track_id"`
	TrackPoints   []*TrackPoint `json:"track_points,omitempty"`
	UpdatedAt     time.Time     `json:"updated_at"`
	UserID        UserID        `json:"user_id"`
	Visibility    int           `json:"visibility"`
	FirstLat      float64       `json:"first_lat"`
	FirstLng      float64       `json:"first_lng"`
	LastLat       float64       `json:"last_lat"`
	LastLng       float64       `json:"last_lng"`
	Metrics       *Metrics      `json:"metrics,omitempty"`
}
