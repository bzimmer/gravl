package cyclinganalytics

import (
	"time"
)

// Fault represents an error response
type Fault struct {
	Message string `json:"error"`
}

func (f *Fault) Error() string {
	return f.Message
}

type (
	UserID   int
	Datetime struct {
		time.Time
	}
)

const (
	// Me represents the user id of the authenticated user
	Me UserID = 0

	// dateTimeFormat used by cyclinganalytics
	datetimeFormat = `"2006-01-02T15:04:05"`
)

func (d *Datetime) UnmarshalJSON(b []byte) (err error) {
	t, err := time.Parse(datetimeFormat, string(b))
	if err != nil {
		return
	}
	d.Time = t
	return
}

func (d *Datetime) MarshalJSON() ([]byte, error) {
	return []byte(d.Time.Format(datetimeFormat)), nil
}

type User struct {
	Email    string `json:"email"`
	ID       UserID `json:"id"`
	Name     string `json:"name"`
	Sex      string `json:"sex"`
	Timezone string `json:"timezone"`
	Units    string `json:"units"`
}

type Has struct {
	Cadence     bool `json:"cadence"`
	Elevation   bool `json:"elevation"`
	GPS         bool `json:"gps"`
	Heartrate   bool `json:"heartrate"`
	Power       bool `json:"power"`
	Speed       bool `json:"speed"`
	Temperature bool `json:"temperature"`
}

type Metadata struct {
	Type string `json:"type"`
}

type Zones struct {
	Heartrate []float64 `json:"heartrate"`
	Power     []float64 `json:"power"`
}

type Summary struct {
	AvgCadence     float64 `json:"avg_cadence"`
	AvgHeartrate   float64 `json:"avg_heartrate"`
	AvgPower       float64 `json:"avg_power"`
	AvgSpeed       float64 `json:"avg_speed"`
	AvgTemperature float64 `json:"avg_temperature"`
	Climbing       float64 `json:"climbing"`
	Decoupling     float64 `json:"decoupling"`
	Distance       float64 `json:"distance"`
	Duration       float64 `json:"duration"`
	Epower         float64 `json:"epower"`
	Intensity      float64 `json:"intensity"`
	Load           float64 `json:"load"`
	LRBalance      float64 `json:"lrbalance"`
	MaxCadence     float64 `json:"max_cadence"`
	MaxHeartrate   float64 `json:"max_heartrate"`
	MaxPower       float64 `json:"max_power"`
	MaxSpeed       float64 `json:"max_speed"`
	MaxTemperature float64 `json:"max_temperature"`
	MinTemperature float64 `json:"min_temperature"`
	MovingTime     float64 `json:"moving_time"`
	PWC150         float64 `json:"pwc150"` // physical working capacity
	PWC170         float64 `json:"pwc170"` // physical working capacity
	PWCR2          float64 `json:"pwc_r2"` // physical working capacity
	RidingTime     float64 `json:"riding_time"`
	TotalTime      float64 `json:"total_time"`
	TRIMP          float64 `json:"trimp"`
	Variability    float64 `json:"variability"`
	Work           float64 `json:"work"`
	Zones          Zones   `json:"zones"`
}

type Shift []int

type Shifts struct {
	Shifts []Shift `json:"shifts"`
}

type Streams struct {
	Power                []float64 `json:"power,omitempty"`
	Speed                []float64 `json:"speed,omitempty"`
	Distance             []float64 `json:"distance,omitempty"`
	Heartrate            []float64 `json:"heartrate,omitempty"`
	Cadence              []float64 `json:"cadence,omitempty"`
	LRBalance            []float64 `json:"lrbalance,omitempty"`
	Latitude             []float64 `json:"latitude,omitempty"`
	Longitude            []float64 `json:"longitude,omitempty"`
	Elevation            []float64 `json:"elevation,omitempty"`
	Gradient             []float64 `json:"gradient,omitempty"`
	Temperature          []float64 `json:"temperature,omitempty"`
	TorqueEffectiveness  []float64 `json:"torque_effectiveness,omitempty"`
	PedalSmoothness      []float64 `json:"pedal_smoothness,omitempty"`
	PlatformCenterOffset []float64 `json:"platform_center_offset,omitempty"`
	PowerPhase           []float64 `json:"power_phase,omitempty"`
	PowerDirection       []float64 `json:"power_direction,omitempty"`
	THB                  []float64 `json:"thb,omitempty"`
	SMO2                 []float64 `json:"smo2,omitempty"`
	RespirationRate      []float64 `json:"respiration_rate,omitempty"`
	HeartRateVariability []float64 `json:"heart_rate_variability,omitempty"`
	Gears                *Shifts   `json:"gears,omitempty"`
}

type Ride struct {
	Format        string   `json:"format"`
	Has           Has      `json:"has"`
	ID            int64    `json:"id"`
	LocalDatetime Datetime `json:"local_datetime"`
	Metadata      Metadata `json:"metadata"`
	Notes         string   `json:"notes"`
	Purpose       string   `json:"purpose"`
	Streams       Streams  `json:"streams"`
	Subtype       string   `json:"subtype"`
	Summary       Summary  `json:"summary"`
	Title         string   `json:"title"`
	Trainer       bool     `json:"trainer"`
	UserID        UserID   `json:"user_id"`
	UTCDatetime   Datetime `json:"utc_datetime"`
}

type Upload struct {
	ID        int64    `json:"upload_id"`
	Status    string   `json:"status"`
	RideID    int64    `json:"ride_id"`
	UserID    UserID   `json:"user_id"`
	Format    string   `json:"format"`
	Datetime  Datetime `json:"datetime"`
	Filename  string   `json:"filename"`
	Size      int64    `json:"size"`
	Error     string   `json:"error"`
	ErrorCode string   `json:"error_code"`
}

type UploadResult struct {
	Upload *Upload `json:"upload"`
	Err    error   `json:"error"`
}
