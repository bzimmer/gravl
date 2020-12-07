package cyclinganalytics

// Fault .
type Fault struct {
	Message string `json:"message"`
}

func (f *Fault) Error() string {
	return f.Message
}

type UserID int

const Me UserID = 0

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
	Lrbalance      float64 `json:"lrbalance"`
	MaxCadence     float64 `json:"max_cadence"`
	MaxHeartrate   float64 `json:"max_heartrate"`
	MaxPower       float64 `json:"max_power"`
	MaxSpeed       float64 `json:"max_speed"`
	MaxTemperature float64 `json:"max_temperature"`
	MinTemperature float64 `json:"min_temperature"`
	MovingTime     float64 `json:"moving_time"`
	Pwc150         float64 `json:"pwc150"`
	Pwc170         float64 `json:"pwc170"`
	PwcR2          float64 `json:"pwc_r2"`
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
	Power                []float64 `json:"power"`
	Speed                []float64 `json:"speed"`
	Distance             []float64 `json:"distance"`
	Heartrate            []float64 `json:"heartrate"`
	Cadence              []float64 `json:"cadence"`
	LRBalance            []float64 `json:"lrbalance"`
	Latitude             []float64 `json:"latitude"`
	Longitude            []float64 `json:"longitude"`
	Elevation            []float64 `json:"elevation"`
	Gradient             []float64 `json:"gradient"`
	Temperature          []float64 `json:"temperature"`
	TorqueEffectiveness  []float64 `json:"torque_effectiveness"`
	PedalSmoothness      []float64 `json:"pedal_smoothness"`
	PlatformCenterOffset []float64 `json:"platform_center_offset"`
	PowerPhase           []float64 `json:"power_phase"`
	PowerDirection       []float64 `json:"power_direction"`
	THB                  []float64 `json:"thb"`
	SMO2                 []float64 `json:"smo2"`
	RespirationRate      []float64 `json:"respiration_rate"`
	HeartRateVariability []float64 `json:"heart_rate_variability"`
	Gears                Shifts    `json:"gears"`
}

type Ride struct {
	Format        string   `json:"format"`
	Has           Has      `json:"has"`
	ID            int64    `json:"id"`
	LocalDatetime string   `json:"local_datetime"`
	Metadata      Metadata `json:"metadata"`
	Notes         string   `json:"notes"`
	Purpose       string   `json:"purpose"`
	Streams       Streams  `json:"streams"`
	Subtype       string   `json:"subtype"`
	Summary       Summary  `json:"summary"`
	Title         string   `json:"title"`
	Trainer       bool     `json:"trainer"`
	UserID        UserID   `json:"user_id"`
	UTCDatetime   string   `json:"utc_datetime"`
}

type RidesResponse struct {
	Rides []*Ride `json:"rides"`
}
