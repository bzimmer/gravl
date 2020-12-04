package rwgps

//go:generate stringer -type=Origin -output=model_string.go

import (
	"fmt"
	"time"

	"github.com/bzimmer/gravl/pkg/common/geo"
)

// Origin of the trip
type Origin int

const (
	// OriginTrip is a ride which was recorded by GPS
	OriginTrip Origin = iota
	// OriginRoute is a ride which was planned on the RWGPS route builder
	OriginRoute
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
	Distance      float64     `json:"distance"`
	Duration      int         `json:"duration"`
	Elevation     *Summary    `json:"ele"`
	ElevationGain float64     `json:"ele_gain"`
	ElevationLoss float64     `json:"ele_loss"`
	EndElevation  float64     `json:"endElevation"`
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
	Watts         interface{} `json:"watts"` // not sure of the real format of Watts
}

// TrackPoint .
type TrackPoint struct {
	Longitude float64 `json:"x"`
	Latitude  float64 `json:"y"`
	Elevation float64 `json:"e"` // elevation in meters
	Distance  float64 `json:"d"` // distance in meters
	Time      float64 `json:"t"` // seconds since epoch, unix timestamp
	Cadence   float64 `json:"c"`
	Grade     float64 `json:"g"`
	Speed     float64 `json:"s"` // kmh ?
}

// Trip .
type Trip struct {
	CreatedAt     time.Time     `json:"created_at"`
	DepartedAt    time.Time     `json:"departed_at"`
	Description   string        `json:"description"`
	Distance      float64       `json:"distance"`
	Duration      int           `json:"duration"`
	ElevationGain float64       `json:"elevation_gain"`
	ElevationLoss float64       `json:"elevation_loss"`
	ID            int64         `json:"id"`
	Name          string        `json:"name"`
	Origin        Origin        `json:"-"`
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

// TripResponse .
type TripResponse struct {
	Type  string `json:"type"`
	Trip  *Trip  `json:"trip"`
	Route *Trip  `json:"route"`
}

type TripsResponse struct {
	Results      []*Trip `json:"results"`
	ResultsCount int     `json:"results_count"`
}

func (t *Trip) Track() (*geo.Track, error) {
	coords := make([][]float64, len(t.TrackPoints))
	for i, tp := range t.TrackPoints {
		coords[i] = []float64{tp.Longitude, tp.Latitude, tp.Elevation}
	}

	var q geo.Origin
	switch t.Origin {
	case OriginTrip:
		q = geo.OriginActivity
	case OriginRoute:
		q = geo.OriginPlanned
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
