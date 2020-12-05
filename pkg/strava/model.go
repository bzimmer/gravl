package strava

import (
	"time"
)

// Error .
type Error struct {
	Resource string `json:"resource"`
	Field    string `json:"field"`
	Code     string `json:"code"`
}

// Fault .
type Fault struct {
	Message string   `json:"message"`
	Errors  []*Error `json:"errors"`
}

func (f *Fault) Error() string {
	return f.Message
}

type StreamMetadata struct {
	OriginalSize int    `json:"original_size"`
	Resolution   string `json:"resolution"`
	SeriesType   string `json:"series_type"`
}

// Stream of data from an activity
type Stream struct {
	StreamMetadata
	Data []float64 `json:"data"`
}

// CoordinateStream of data from an activity
type CoordinateStream struct {
	StreamMetadata
	Data [][]float64 `json:"data"`
}

type Streams struct {
	ActivityID  int64             `json:"activity_id"`
	LatLng      *CoordinateStream `json:"latlng"`
	Altitude    *Stream           `json:"altitude"`
	Time        *Stream           `json:"time"`
	Distance    *Stream           `json:"distance"`
	Velocity    *Stream           `json:"velocity_smooth"`
	HeartRate   *Stream           `json:"heartrate"`
	Cadence     *Stream           `json:"cadence"`
	Watts       *Stream           `json:"watts"`
	Temperature *Stream           `json:"temp"`
	Moving      *Stream           `json:"moving"`
	Grade       *Stream           `json:"grade_smooth"`
}

// Gear represents gear used by the athlete
type Gear struct {
	ID            string  `json:"id"`
	Primary       bool    `json:"primary"`
	Name          string  `json:"name"`
	ResourceState int     `json:"resource_state"`
	Distance      float64 `json:"distance"`
	AthleteID     int     `json:"athlete_id"`
}

// Totals .
type Totals struct {
	Distance         float64 `json:"distance"`
	AchievementCount int     `json:"achievement_count"`
	Count            int     `json:"count"`
	ElapsedTime      int     `json:"elapsed_time"`
	ElevationGain    float64 `json:"elevation_gain"`
	MovingTime       int     `json:"moving_time"`
}

// Stats .
type Stats struct {
	RecentRunTotals           *Totals `json:"recent_run_totals"`
	AllRunTotals              *Totals `json:"all_run_totals"`
	RecentSwimTotals          *Totals `json:"recent_swim_totals"`
	BiggestRideDistance       float64 `json:"biggest_ride_distance"`
	YtdSwimTotals             *Totals `json:"ytd_swim_totals"`
	AllSwimTotals             *Totals `json:"all_swim_totals"`
	RecentRideTotals          *Totals `json:"recent_ride_totals"`
	BiggestClimbElevationGain float64 `json:"biggest_climb_elevation_gain"`
	YtdRideTotals             *Totals `json:"ytd_ride_totals"`
	AllRideTotals             *Totals `json:"all_ride_totals"`
	YtdRunTotals              *Totals `json:"ytd_run_totals"`
}

// Athlete represents a Strava athlete
type Athlete struct {
	ID                    int           `json:"id"`
	Username              string        `json:"username"`
	ResourceState         int           `json:"resource_state"`
	Firstname             string        `json:"firstname"`
	Lastname              string        `json:"lastname"`
	City                  string        `json:"city"`
	State                 string        `json:"state"`
	Country               string        `json:"country"`
	Sex                   string        `json:"sex"`
	Premium               bool          `json:"premium"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
	BadgeTypeID           int           `json:"badge_type_id"`
	ProfileMedium         string        `json:"profile_medium"`
	Profile               string        `json:"profile"`
	Friend                interface{}   `json:"friend"`
	Follower              interface{}   `json:"follower"`
	FollowerCount         int           `json:"follower_count"`
	FriendCount           int           `json:"friend_count"`
	MutualFriendCount     int           `json:"mutual_friend_count"`
	AthleteType           int           `json:"athlete_type"`
	DatePreference        string        `json:"date_preference"`
	MeasurementPreference string        `json:"measurement_preference"`
	Clubs                 []interface{} `json:"clubs"`
	Ftp                   float64       `json:"ftp"`
	Weight                float64       `json:"weight"`
	Bikes                 []*Gear       `json:"bikes"`
	Shoes                 []*Gear       `json:"shoes"`
}

// Map .
type Map struct {
	ID              string `json:"id"`
	Polyline        string `json:"polyline"`
	ResourceState   int    `json:"resource_state"`
	SummaryPolyline string `json:"summary_polyline"`
}

// Lap .
type Lap struct {
	ID                 int64     `json:"id"`
	ResourceState      int       `json:"resource_state"`
	Name               string    `json:"name"`
	Activity           *Activity `json:"activity"`
	Athlete            *Athlete  `json:"athlete"`
	ElapsedTime        int       `json:"elapsed_time"`
	MovingTime         int       `json:"moving_time"`
	StartDate          time.Time `json:"start_date"`
	StartDateLocal     time.Time `json:"start_date_local"`
	Distance           float64   `json:"distance"`
	StartIndex         int       `json:"start_index"`
	EndIndex           int       `json:"end_index"`
	TotalElevationGain float64   `json:"total_elevation_gain"`
	AverageSpeed       float64   `json:"average_speed"`
	MaxSpeed           float64   `json:"max_speed"`
	AverageCadence     float64   `json:"average_cadence"`
	DeviceWatts        bool      `json:"device_watts"`
	AverageWatts       float64   `json:"average_watts"`
	LapIndex           int       `json:"lap_index"`
	Split              int       `json:"split"`
}

// Segment .
type Segment struct {
	ID                 int       `json:"id"`
	ResourceState      int       `json:"resource_state"`
	Name               string    `json:"name"`
	ActivityType       string    `json:"activity_type"`
	Distance           float64   `json:"distance"`
	AverageGrade       float64   `json:"average_grade"`
	MaximumGrade       float64   `json:"maximum_grade"`
	ElevationHigh      float64   `json:"elevation_high"`
	ElevationLow       float64   `json:"elevation_low"`
	StartLatlng        []float64 `json:"start_latlng"`
	EndLatlng          []float64 `json:"end_latlng"`
	ClimbCategory      int       `json:"climb_category"`
	City               string    `json:"city"`
	State              string    `json:"state"`
	Country            string    `json:"country"`
	Private            bool      `json:"private"`
	Hazardous          bool      `json:"hazardous"`
	Starred            bool      `json:"starred"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	TotalElevationGain int       `json:"total_elevation_gain"`
	Map                *Map      `json:"map"`
	EffortCount        int       `json:"effort_count"`
	AthleteCount       int       `json:"athlete_count"`
	StarCount          int       `json:"star_count"`
	AthletePrEffort    struct {
		Distance       float64   `json:"distance"`
		StartDateLocal time.Time `json:"start_date_local"`
		ActivityID     int       `json:"activity_id"`
		ElapsedTime    int       `json:"elapsed_time"`
		IsKom          bool      `json:"is_kom"`
		ID             int       `json:"id"`
		StartDate      time.Time `json:"start_date"`
	} `json:"athlete_pr_effort"`
	AthleteSegmentStats struct {
		PrElapsedTime int       `json:"pr_elapsed_time"`
		PrDate        time.Time `json:"pr_date"`
		EffortCount   int       `json:"effort_count"`
		PrActivityID  int       `json:"pr_activity_id"`
	} `json:"athlete_segment_stats"`
}

// MetaActivity .
type MetaActivity struct {
	ID            int64 `json:"id"`
	ResourceState int   `json:"resource_state"`
}

// SegmentEffort .
type SegmentEffort struct {
	ID             int64         `json:"id"`
	ResourceState  int           `json:"resource_state"`
	Name           string        `json:"name"`
	Activity       *MetaActivity `json:"activity"`
	Athlete        *Athlete      `json:"athlete"`
	ElapsedTime    int           `json:"elapsed_time"`
	MovingTime     int           `json:"moving_time"`
	StartDate      time.Time     `json:"start_date"`
	StartDateLocal time.Time     `json:"start_date_local"`
	Distance       float64       `json:"distance"`
	StartIndex     int           `json:"start_index"`
	EndIndex       int           `json:"end_index"`
	AverageCadence float64       `json:"average_cadence"`
	DeviceWatts    bool          `json:"device_watts"`
	AverageWatts   float64       `json:"average_watts"`
	Segment        *Segment      `json:"segment"`
	KomRank        int           `json:"kom_rank"`
	PrRank         int           `json:"pr_rank"`
	Achievements   []interface{} `json:"achievements"`
	Hidden         bool          `json:"hidden"`
}

// SplitsMetric .
type SplitsMetric struct {
	Distance            float64 `json:"distance"`
	ElapsedTime         int     `json:"elapsed_time"`
	ElevationDifference float64 `json:"elevation_difference"`
	MovingTime          int     `json:"moving_time"`
	Split               int     `json:"split"`
	AverageSpeed        float64 `json:"average_speed"`
	PaceZone            int     `json:"pace_zone"`
}

// HighlightedKudosers .
type HighlightedKudosers struct {
	DestinationURL string `json:"destination_url"`
	DisplayName    string `json:"display_name"`
	AvatarURL      string `json:"avatar_url"`
	ShowName       bool   `json:"show_name"`
}

// Photos .
type Photos struct {
	Primary struct {
		ID       interface{} `json:"id"`
		UniqueID string      `json:"unique_id"`
		Urls     struct {
			Num100 string `json:"100"`
			Num600 string `json:"600"`
		} `json:"urls"`
		Source int `json:"source"`
	} `json:"primary"`
	UsePrimaryPhoto bool `json:"use_primary_photo"`
	Count           int  `json:"count"`
}

// Activity represents an activity
type Activity struct {
	ID                       int64                  `json:"id"`
	ResourceState            int                    `json:"resource_state"`
	ExternalID               string                 `json:"external_id"`
	UploadID                 int64                  `json:"upload_id"`
	Athlete                  *Athlete               `json:"athlete"`
	Name                     string                 `json:"name"`
	Distance                 float64                `json:"distance"`
	MovingTime               int                    `json:"moving_time"`
	ElapsedTime              int                    `json:"elapsed_time"`
	TotalElevationGain       float64                `json:"total_elevation_gain"`
	Type                     string                 `json:"type"`
	StartDate                time.Time              `json:"start_date"`
	StartDateLocal           time.Time              `json:"start_date_local"`
	Timezone                 string                 `json:"timezone"`
	UtcOffset                float64                `json:"utc_offset"`
	StartLatlng              []float64              `json:"start_latlng"`
	EndLatlng                []float64              `json:"end_latlng"`
	LocationCity             string                 `json:"location_city"`
	LocationState            string                 `json:"location_state"`
	LocationCountry          string                 `json:"location_country"`
	AchievementCount         int                    `json:"achievement_count"`
	KudosCount               int                    `json:"kudos_count"`
	CommentCount             int                    `json:"comment_count"`
	AthleteCount             int                    `json:"athlete_count"`
	PhotoCount               int                    `json:"photo_count"`
	Map                      *Map                   `json:"map"`
	Trainer                  bool                   `json:"trainer"`
	Commute                  bool                   `json:"commute"`
	Manual                   bool                   `json:"manual"`
	Private                  bool                   `json:"private"`
	Flagged                  bool                   `json:"flagged"`
	GearID                   string                 `json:"gear_id"`
	FromAcceptedTag          bool                   `json:"from_accepted_tag"`
	AverageSpeed             float64                `json:"average_speed"`
	MaxSpeed                 float64                `json:"max_speed"`
	AverageCadence           float64                `json:"average_cadence"`
	AverageTemp              int                    `json:"average_temp"`
	AverageWatts             float64                `json:"average_watts"`
	WeightedAverageWatts     int                    `json:"weighted_average_watts"`
	Kilojoules               float64                `json:"kilojoules"`
	DeviceWatts              bool                   `json:"device_watts"`
	HasHeartrate             bool                   `json:"has_heartrate"`
	MaxWatts                 int                    `json:"max_watts"`
	ElevHigh                 float64                `json:"elev_high"`
	ElevLow                  float64                `json:"elev_low"`
	PrCount                  int                    `json:"pr_count"`
	TotalPhotoCount          int                    `json:"total_photo_count"`
	HasKudoed                bool                   `json:"has_kudoed"`
	WorkoutType              int                    `json:"workout_type"`
	SufferScore              float64                `json:"suffer_score"`
	Description              string                 `json:"description"`
	Calories                 float64                `json:"calories"`
	SegmentEfforts           []*SegmentEffort       `json:"segment_efforts"`
	SplitsMetric             []*SplitsMetric        `json:"splits_metric"`
	Laps                     []*Lap                 `json:"laps"`
	Gear                     *Gear                  `json:"gear"`
	PartnerBrandTag          interface{}            `json:"partner_brand_tag"`
	Photos                   *Photos                `json:"photos"`
	HighlightedKudosers      []*HighlightedKudosers `json:"highlighted_kudosers"`
	DeviceName               string                 `json:"device_name"`
	EmbedToken               string                 `json:"embed_token"`
	SegmentLeaderboardOptOut bool                   `json:"segment_leaderboard_opt_out"`
	LeaderboardOptOut        bool                   `json:"leaderboard_opt_out"`
	PerceivedExertion        float64                `json:"perceived_exertion"`
	PreferPerceivedExertion  bool                   `json:"prefer_perceived_exertion"`
}

// Route .
type Route struct {
	Private             bool       `json:"private"`
	Distance            float64    `json:"distance"`
	Athlete             *Athlete   `json:"athlete"`
	Description         string     `json:"description"`
	CreatedAt           time.Time  `json:"created_at"`
	ElevationGain       float64    `json:"elevation_gain"`
	Type                int        `json:"type"`
	EstimatedMovingTime int        `json:"estimated_moving_time"`
	Segments            []*Segment `json:"segments"`
	Starred             bool       `json:"starred"`
	UpdatedAt           time.Time  `json:"updated_at"`
	SubType             int        `json:"sub_type"`
	IDStr               string     `json:"id_str"`
	Name                string     `json:"name"`
	ID                  int        `json:"id"`
	Map                 *Map       `json:"map"`
	Timestamp           int        `json:"timestamp"`
}
