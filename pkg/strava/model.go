package strava

//go:generate genny -in ../common/fp/template/every.go -out "strava_fp_every.go" -pkg "strava" gen "ValueType=Activity"
//go:generate genny -in ../common/fp/template/filter.go -out "strava_fp_filter.go" -pkg "strava" gen "ValueType=Activity"
//go:generate genny -in ../common/fp/template/groupby.go -out "strava_fp_groupby.go" -pkg "strava" gen "ValueType=Activity KeyType=int"
//go:generate genny -in ../common/fp/template/map.go -out "strava_fp_map.go" -pkg "strava" gen "ValueType=Activity"
//go:generate genny -in ../common/fp/template/reduce.go -out "strava_fp_reduce.go" -pkg "strava" gen "ValueType=Activity"

import (
	"time"

	"github.com/martinlindhe/unit"
	"github.com/twpayne/go-geom"
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

// SpeedStream of data from an activity
type SpeedStream struct {
	StreamMetadata
	Data []unit.Speed `json:"data" units:"kph"`
}

// LengthStream of data from an activity
type LengthStream struct {
	StreamMetadata
	Data []unit.Length `json:"data" units:"m"`
}

// BoolStream of data from an activity
type BoolStream struct {
	StreamMetadata
	Data []bool `json:"data"`
}

type Streams struct {
	ActivityID  int64             `json:"activity_id"`
	LatLng      *CoordinateStream `json:"latlng,omitempty"`
	Elevation   *LengthStream     `json:"altitude,omitempty"`
	Time        *Stream           `json:"time,omitempty"`
	Distance    *LengthStream     `json:"distance,omitempty"`
	Velocity    *SpeedStream      `json:"velocity_smooth,omitempty"`
	HeartRate   *Stream           `json:"heartrate,omitempty"`
	Cadence     *Stream           `json:"cadence,omitempty"`
	Watts       *Stream           `json:"watts,omitempty"`
	Temperature *Stream           `json:"temp,omitempty"`
	Moving      *BoolStream       `json:"moving,omitempty"`
	Grade       *Stream           `json:"grade_smooth,omitempty"`
}

// Gear represents gear used by the athlete
type Gear struct {
	ID            string      `json:"id"`
	Primary       bool        `json:"primary"`
	Name          string      `json:"name"`
	ResourceState int         `json:"resource_state"`
	Distance      unit.Length `json:"distance" units:"m"`
	AthleteID     int         `json:"athlete_id"`
}

// Totals .
type Totals struct {
	Distance         unit.Length `json:"distance" units:"m"`
	AchievementCount int         `json:"achievement_count"`
	Count            int         `json:"count"`
	ElapsedTime      int         `json:"elapsed_time" units:"s"`
	ElevationGain    unit.Length `json:"elevation_gain" units:"m"`
	MovingTime       int         `json:"moving_time"`
}

// Stats .
type Stats struct {
	RecentRunTotals           *Totals     `json:"recent_run_totals"`
	AllRunTotals              *Totals     `json:"all_run_totals"`
	RecentSwimTotals          *Totals     `json:"recent_swim_totals"`
	BiggestRideDistance       unit.Length `json:"biggest_ride_distance" units:"m"`
	YtdSwimTotals             *Totals     `json:"ytd_swim_totals"`
	AllSwimTotals             *Totals     `json:"all_swim_totals"`
	RecentRideTotals          *Totals     `json:"recent_ride_totals"`
	BiggestClimbElevationGain unit.Length `json:"biggest_climb_elevation_gain" units:"m"`
	YtdRideTotals             *Totals     `json:"ytd_ride_totals"`
	AllRideTotals             *Totals     `json:"all_ride_totals"`
	YtdRunTotals              *Totals     `json:"ytd_run_totals"`
}

type Club struct {
	Admin           bool   `json:"admin"`
	City            string `json:"city"`
	Country         string `json:"country"`
	CoverPhoto      string `json:"cover_photo"`
	CoverPhotoSmall string `json:"cover_photo_small"`
	Featured        bool   `json:"featured"`
	ID              int    `json:"id"`
	MemberCount     int    `json:"member_count"`
	Membership      string `json:"membership"`
	Name            string `json:"name"`
	Owner           bool   `json:"owner"`
	Private         bool   `json:"private"`
	Profile         string `json:"profile"`
	ProfileMedium   string `json:"profile_medium"`
	ResourceState   int    `json:"resource_state"`
	SportType       string `json:"sport_type"`
	State           string `json:"state"`
	URL             string `json:"url"`
	Verified        bool   `json:"verified"`
}

// Athlete represents a Strava athlete
type Athlete struct {
	ID                    int         `json:"id"`
	Username              string      `json:"username"`
	ResourceState         int         `json:"resource_state"`
	Firstname             string      `json:"firstname"`
	Lastname              string      `json:"lastname"`
	City                  string      `json:"city"`
	State                 string      `json:"state"`
	Country               string      `json:"country"`
	Sex                   string      `json:"sex"`
	Premium               bool        `json:"premium"`
	CreatedAt             time.Time   `json:"created_at"`
	UpdatedAt             time.Time   `json:"updated_at"`
	BadgeTypeID           int         `json:"badge_type_id"`
	ProfileMedium         string      `json:"profile_medium"`
	Profile               string      `json:"profile"`
	Friend                interface{} `json:"friend"`
	Follower              interface{} `json:"follower"`
	FollowerCount         int         `json:"follower_count"`
	FriendCount           int         `json:"friend_count"`
	MutualFriendCount     int         `json:"mutual_friend_count"`
	AthleteType           int         `json:"athlete_type"`
	DatePreference        string      `json:"date_preference"`
	MeasurementPreference string      `json:"measurement_preference"`
	Clubs                 []*Club     `json:"clubs"`
	FTP                   float64     `json:"ftp"`
	Weight                float64     `json:"weight" units:"kg"`
	Bikes                 []*Gear     `json:"bikes"`
	Shoes                 []*Gear     `json:"shoes"`
}

// Map .
type Map struct {
	ID              string `json:"id"`
	Polyline        string `json:"polyline"`
	ResourceState   int    `json:"resource_state"`
	SummaryPolyline string `json:"summary_polyline"`
}

func (m *Map) LineString() (*geom.LineString, error) {
	return polylineToLineString(m.Polyline, m.SummaryPolyline)
}

// Lap .
type Lap struct {
	ID                 int64       `json:"id"`
	ResourceState      int         `json:"resource_state"`
	Name               string      `json:"name"`
	Activity           *Activity   `json:"activity"`
	Athlete            *Athlete    `json:"athlete"`
	ElapsedTime        int         `json:"elapsed_time"`
	MovingTime         int         `json:"moving_time"`
	StartDate          time.Time   `json:"start_date"`
	StartDateLocal     time.Time   `json:"start_date_local"`
	Distance           unit.Length `json:"distance" units:"m"`
	StartIndex         int         `json:"start_index"`
	EndIndex           int         `json:"end_index"`
	TotalElevationGain unit.Length `json:"total_elevation_gain" units:"m"`
	AverageSpeed       unit.Speed  `json:"average_speed" units:"kph"`
	MaxSpeed           unit.Speed  `json:"max_speed" units:"kph"`
	AverageCadence     float64     `json:"average_cadence"`
	DeviceWatts        bool        `json:"device_watts"`
	AverageWatts       float64     `json:"average_watts"`
	LapIndex           int         `json:"lap_index"`
	Split              int         `json:"split"`
}

type PREffort struct {
	Distance       unit.Length `json:"distance" units:"m"`
	StartDateLocal time.Time   `json:"start_date_local"`
	ActivityID     int         `json:"activity_id"`
	ElapsedTime    int         `json:"elapsed_time"`
	IsKOM          bool        `json:"is_kom"`
	ID             int         `json:"id"`
	StartDate      time.Time   `json:"start_date"`
}

type SegmentStats struct {
	PRElapsedTime int       `json:"pr_elapsed_time"`
	PRDate        time.Time `json:"pr_date"`
	EffortCount   int       `json:"effort_count"`
	PRActivityID  int       `json:"pr_activity_id"`
}

// Segment .
type Segment struct {
	ID                  int           `json:"id"`
	ResourceState       int           `json:"resource_state"`
	Name                string        `json:"name"`
	ActivityType        string        `json:"activity_type"`
	Distance            unit.Length   `json:"distance" units:"m"`
	AverageGrade        float64       `json:"average_grade"`
	MaximumGrade        float64       `json:"maximum_grade"`
	ElevationHigh       unit.Length   `json:"elevation_high" units:"m"`
	ElevationLow        unit.Length   `json:"elevation_low" units:"m"`
	StartLatlng         []float64     `json:"start_latlng"`
	EndLatlng           []float64     `json:"end_latlng"`
	ClimbCategory       int           `json:"climb_category"`
	City                string        `json:"city"`
	State               string        `json:"state"`
	Country             string        `json:"country"`
	Private             bool          `json:"private"`
	Hazardous           bool          `json:"hazardous"`
	Starred             bool          `json:"starred"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
	ElevationGain       unit.Length   `json:"total_elevation_gain" units:"m"`
	Map                 *Map          `json:"map"`
	EffortCount         int           `json:"effort_count"`
	AthleteCount        int           `json:"athlete_count"`
	StarCount           int           `json:"star_count"`
	PREffort            *PREffort     `json:"athlete_pr_effort"`
	AthleteSegmentStats *SegmentStats `json:"athlete_segment_stats"`
}

// MetaActivity .
type MetaActivity struct {
	ID            int64 `json:"id"`
	ResourceState int   `json:"resource_state"`
}

// Achievement .
type Achievement struct {
	Rank   int    `json:"rank"`
	Type   string `json:"type"`
	TypeID int    `json:"type_id"`
}

// SegmentEffort .
type SegmentEffort struct {
	ID             int64          `json:"id"`
	ResourceState  int            `json:"resource_state"`
	Name           string         `json:"name"`
	Activity       *MetaActivity  `json:"activity"`
	Athlete        *Athlete       `json:"athlete"`
	ElapsedTime    int            `json:"elapsed_time"`
	MovingTime     int            `json:"moving_time"`
	StartDate      time.Time      `json:"start_date"`
	StartDateLocal time.Time      `json:"start_date_local"`
	Distance       unit.Length    `json:"distance" units:"m"`
	StartIndex     int            `json:"start_index"`
	EndIndex       int            `json:"end_index"`
	AverageCadence float64        `json:"average_cadence"`
	DeviceWatts    bool           `json:"device_watts"`
	AverageWatts   float64        `json:"average_watts"`
	Segment        *Segment       `json:"segment"`
	KOMRank        int            `json:"kom_rank"`
	PRRank         int            `json:"pr_rank"`
	Achievements   []*Achievement `json:"achievements"`
	Hidden         bool           `json:"hidden"`
}

// SplitsMetric .
type SplitsMetric struct {
	Distance            unit.Length `json:"distance" units:"m"`
	ElapsedTime         int         `json:"elapsed_time"`
	ElevationDifference unit.Length `json:"elevation_difference" units:"m"`
	MovingTime          int         `json:"moving_time"`
	Split               int         `json:"split"`
	AverageSpeed        float64     `json:"average_speed"`
	PaceZone            int         `json:"pace_zone"`
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
	Distance                 unit.Length            `json:"distance" units:"m"`
	MovingTime               int                    `json:"moving_time"`
	ElapsedTime              int                    `json:"elapsed_time"`
	ElevationGain            unit.Length            `json:"total_elevation_gain" units:"m"`
	Type                     string                 `json:"type"`
	StartDate                time.Time              `json:"start_date"`
	StartDateLocal           time.Time              `json:"start_date_local"`
	Timezone                 string                 `json:"timezone"`
	UTCOffset                float64                `json:"utc_offset"`
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
	AverageSpeed             unit.Speed             `json:"average_speed" units:"kph"`
	MaxSpeed                 unit.Speed             `json:"max_speed" units:"kph"`
	AverageCadence           float64                `json:"average_cadence"`
	AverageTemperature       unit.Temperature       `json:"average_temp" units:"C"`
	AverageWatts             float64                `json:"average_watts"`
	WeightedAverageWatts     int                    `json:"weighted_average_watts"`
	Kilojoules               float64                `json:"kilojoules"`
	DeviceWatts              bool                   `json:"device_watts"`
	HasHeartrate             bool                   `json:"has_heartrate"`
	MaxWatts                 int                    `json:"max_watts"`
	ElevationHigh            unit.Length            `json:"elev_high" units:"m"`
	ElevationLow             unit.Length            `json:"elev_low" units:"m"`
	PRCount                  int                    `json:"pr_count"`
	TotalPhotoCount          int                    `json:"total_photo_count"`
	HasKudoed                bool                   `json:"has_kudoed"`
	WorkoutType              int                    `json:"workout_type"`
	SufferScore              float64                `json:"suffer_score"`
	Description              string                 `json:"description"`
	Calories                 float64                `json:"calories"`
	SegmentEfforts           []*SegmentEffort       `json:"segment_efforts,omitempty"`
	SplitsMetric             []*SplitsMetric        `json:"splits_metric,omitempty"`
	Laps                     []*Lap                 `json:"laps,omitempty"`
	Gear                     *Gear                  `json:"gear,omitempty"`
	PartnerBrandTag          interface{}            `json:"partner_brand_tag"`
	Photos                   *Photos                `json:"photos,omitempty"`
	HighlightedKudosers      []*HighlightedKudosers `json:"highlighted_kudosers,omitempty"`
	DeviceName               string                 `json:"device_name"`
	EmbedToken               string                 `json:"embed_token"`
	SegmentLeaderboardOptOut bool                   `json:"segment_leaderboard_opt_out"`
	LeaderboardOptOut        bool                   `json:"leaderboard_opt_out"`
	PerceivedExertion        float64                `json:"perceived_exertion"`
	PreferPerceivedExertion  bool                   `json:"prefer_perceived_exertion"`
	Streams                  *Streams               `json:"streams,omitempty"`
}

// Route .
type Route struct {
	Private             bool        `json:"private"`
	Distance            unit.Length `json:"distance" units:"m"`
	Athlete             *Athlete    `json:"athlete"`
	Description         string      `json:"description"`
	CreatedAt           time.Time   `json:"created_at"`
	ElevationGain       unit.Length `json:"elevation_gain" units:"m"`
	Type                int         `json:"type"`
	EstimatedMovingTime int         `json:"estimated_moving_time"`
	Segments            []*Segment  `json:"segments"`
	Starred             bool        `json:"starred"`
	UpdatedAt           time.Time   `json:"updated_at"`
	SubType             int         `json:"sub_type"`
	IDStr               string      `json:"id_str"`
	Name                string      `json:"name"`
	ID                  int         `json:"id"`
	Map                 *Map        `json:"map"`
	Timestamp           int         `json:"timestamp"`
}
