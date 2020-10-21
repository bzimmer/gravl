package rwgps

// UserResponse .
type UserResponse struct {
	AdditionalDrawerItems *AdditionalDrawerItems `json:"additional_drawer_items"`
	Labs                  *Labs                  `json:"labs"`
	User                  *User                  `json:"user"`
}

// AdditionalDrawerItems .
type AdditionalDrawerItems []struct {
	Icon   string `json:"icon"`
	Name   string `json:"name"`
	Target string `json:"target"`
}

// Labs .
type Labs struct {
	CanSeeHomepageOptIn bool `json:"can_see_homepage_opt_in"`
	FitExport           bool `json:"fit_export"`
	HomepageRedesign    bool `json:"homepage_redesign"`
	MobileHomeRedesign  bool `json:"mobile_home_redesign"`
	PinnedInSearchbar   bool `json:"pinned_in_searchbar"`
	ReliveSync          bool `json:"relive_sync"`
	RoutePlanner        bool `json:"route_planner"`
	RspRedesign         bool `json:"rsp_redesign"`
	SmartExport         bool `json:"smart_export"`
	StravaSync          bool `json:"strava_sync"`
	TrspRedesign        bool `json:"trsp_redesign"`
	Varia               bool `json:"varia"`
}

// Gear .
type Gear struct {
	Archived          bool   `json:"archived"`
	CreatedAt         string `json:"created_at"`
	Description       string `json:"description"`
	ExcludeFromTotals bool   `json:"exclude_from_totals"`
	GearModelID       int64  `json:"gear_model_id"`
	GearTypeID        string `json:"gear_type_id"`
	GroupMembershipID int64  `json:"group_membership_id"`
	ID                int64  `json:"id"`
	Make              string `json:"make"`
	Model             string `json:"model"`
	Name              string `json:"name"`
	Nickname          string `json:"nickname"`
	// Project529ID      interface{} `json:"project_529_id"`
	SerialNumber    string  `json:"serial_number"`
	UseUserTimezone bool    `json:"use_user_timezone"`
	Visibility      int64   `json:"visibility"`
	Weight          float64 `json:"weight"`
	Year            int64   `json:"year"`
}

// User .
type User struct {
	AccountLevel       int64   `json:"account_level"`
	AdministrativeArea string  `json:"administrative_area"`
	Age                int     `json:"age"`
	AuthToken          string  `json:"auth_token"`
	ClubIds            []int64 `json:"club_ids"`
	CountryCode        string  `json:"country_code"`
	CreatedAt          string  `json:"created_at"`
	// Deactivated        *int    `json:"deactivated"`
	// DeactivatedAt           interface{} `json:"deactivated_at"`
	Description string `json:"description"`
	DisplayName string `json:"display_name"`
	// Dob                     interface{} `json:"dob"`
	// EligibleForOnetimeTrial bool   `json:"eligible_for_onetime_trial"`
	Email string `json:"email"`
	// EmailBounceCount        int64  `json:"email_bounce_count"`
	// EmailOnComment          bool   `json:"email_on_comment"`
	// EmailOnMessage          bool   `json:"email_on_message"`
	// EmailOnUpdate           bool   `json:"email_on_update"`
	// EmailVisible            bool   `json:"email_visible"`
	// FacebookID              interface{} `json:"facebook_id"`
	FirstName string  `json:"first_name"`
	Gear      *[]Gear `json:"gear"`
	// HasEvents bool    `json:"has_events"`
	// HeatmapOptOut      bool    `json:"heatmap_optout"`
	// HighlightedPhotoID int64 `json:"highlighted_photo_id"`
	// HrMax              interface{} `json:"hr_max"`
	// HrRest             interface{} `json:"hr_rest"`
	// HrZone1High        interface{} `json:"hr_zone_1_high"`
	// HrZone1Low         interface{} `json:"hr_zone_1_low"`
	// HrZone2High        interface{} `json:"hr_zone_2_high"`
	// HrZone2Low         interface{} `json:"hr_zone_2_low"`
	// HrZone3High        interface{} `json:"hr_zone_3_high"`
	// HrZone3Low         interface{} `json:"hr_zone_3_low"`
	// HrZone4High        interface{} `json:"hr_zone_4_high"`
	// HrZone4Low         interface{} `json:"hr_zone_4_low"`
	// HrZone5High        interface{} `json:"hr_zone_5_high"`
	// HrZone5Low         interface{} `json:"hr_zone_5_low"`
	ID int64 `json:"id"`
	// Interests          string  `json:"interests"`
	// IsMale             bool    `json:"is_male"`
	// IsShadowUser       bool    `json:"is_shadow_user"`
	// LastLoginAt        string  `json:"last_login_at"`
	LastName    string  `json:"last_name"`
	Latitude    float64 `json:"lat"`
	Longitude   float64 `json:"lng"`
	Locale      string  `json:"locale"`
	Locality    string  `json:"locality"`
	MetricUnits bool    `json:"metric_units"`
	Name        string  `json:"name"`
	// NeedsPasswordReset bool    `json:"needs_password_reset"`
	// NumUnreadMessages  int64   `json:"num_unread_messages"`
	PostalCode string `json:"postal_code"`
	// Preferences             struct {
	// 	SkipInstallPrompt     bool   `json:"S2D:skipInstallPrompt"`
	// 	CalendarShowHr        bool   `json:"calendar_show_hr"`
	// 	CollapseUserDashboard bool   `json:"collapse_user_dashboard"`
	// 	DashboardHome         string `json:"dashboard_home"`
	// 	DateFormat            string `json:"date_format"`
	// 	DefaultCareerInterval string `json:"default_career_interval"`
	DefaultGearID string `json:"default_gear_id"`
	// 	DefaultPrivacyRoute   int64  `json:"default_privacy_route"`
	// 	DefaultPrivacyTrip    int64  `json:"default_privacy_trip"`
	// 	Facebook              struct {
	// 		NotifyOnActivity bool `json:"notify_on_activity"`
	// 		NotifyOnRoute    bool `json:"notify_on_route"`
	// 	} `json:"facebook"`
	// 	FeedsDefaultTab                            string  `json:"feeds_default_tab"`
	// 	GlobalNotificationDismissed                int64   `json:"global_notification_dismissed"`
	// 	HideZendeskButton                          bool    `json:"hide_zendesk_button"`
	// 	Lng                                        float64 `json:"lng"`
	// 	MetricUnits                                bool    `json:"metric_units"`
	// 	PlannerOnboardingDismissed                 bool    `json:"planner-onboarding-dismissed"`
	// 	PlannerOverlay                             string  `json:"planner_overlay"`
	// 	PrintIncludePOI                            bool    `json:"print.include_poi"`
	// 	ProfileDefaultTab                          string  `json:"profile_default_tab"`
	// 	ProfileSelectedTab                         string  `json:"profile_selected_tab"`
	// 	ReceiveSegmentNotifications                bool    `json:"receive_segment_notifications"`
	// 	RouteExportAdvancedTurnNotificationChecked bool    `json:"route_export_advanced_turn_notification_checked"`
	// 	RoutePlannerDirectionsType                 string  `json:"route_planner_directions_type"`
	// 	RoutePlannerLeftSidebarClosed              bool    `json:"route_planner_left_sidebar_closed"`
	// 	RoutePlannerRightSidebarClosed             bool    `json:"route_planner_right_sidebar_closed"`
	// 	RouteViewerActiveSubtab                    string  `json:"route_viewer_active_subtab"`
	// 	RouteViewerActiveTab                       string  `json:"route_viewer_active_tab"`
	// 	RouteViewerEnableDistanceMarkers           bool    `json:"route_viewer_enable_distance_markers"`
	// 	RoutesGridHeight                           string  `json:"routes_grid_height"`
	// 	RROnboardingDismissed                      bool    `json:"rr-onboarding-dismissed"`
	// 	SegmentsPrivate                            bool    `json:"segments_private"`
	// 	ShowDashExpiringTrialNotice                bool    `json:"show_dash_expiring_trial_notice"`
	// 	ShowDashMiniOnboarding                     bool    `json:"show_dash_mini_onboarding"`
	// 	ShowDashOnboardingOverlay                  bool    `json:"show_dash_onboarding_overlay"`
	// 	ShowDashTrialNotice                        bool    `json:"show_dash_trial_notice"`
	// 	SmartExportFileFormat                      string  `json:"smart_export_file_format"`
	// 	SmartExportState                           string  `json:"smart_export_state"`
	// 	TrspNoticeDismissed                        bool    `json:"trsp_notice_dismissed"`
	// } `json:"preferences"`
	// Privileges []string `json:"privileges"`
	// PushApplications []interface{} `json:"push_applications"`
	// RelevantGoalParticipants []struct {
	// 	AmountCompleted int64 `json:"amount_completed"`
	// 	Goal            struct {
	// 		Description string `json:"description"`
	// 		Icon        string `json:"icon"`
	// 		IconSmall   string `json:"icon_small"`
	// 		ID          int64  `json:"id"`
	// 		Name        string `json:"name"`
	// 		User        struct {
	// 			AccountLevel             int64       `json:"account_level"`
	// 			AdministrativeArea       string      `json:"administrative_area"`
	// 			Age                      interface{} `json:"age"`
	// 			CreatedAt                string      `json:"created_at"`
	// 			HighlightedPhotoChecksum string      `json:"highlighted_photo_checksum"`
	// 			HighlightedPhotoID       int64       `json:"highlighted_photo_id"`
	// 			ID                       int64       `json:"id"`
	// 			Lat                      float64     `json:"lat"`
	// 			Lng                      float64     `json:"lng"`
	// 			Locality                 string      `json:"locality"`
	// 			Name                     string      `json:"name"`
	// 		} `json:"user"`
	// 	} `json:"goal"`
	// 	GoalParams struct {
	// 		Percent      float64 `json:"percent"`
	// 		ProgressText string  `json:"progress_text"`
	// 		Trailer      string  `json:"trailer"`
	// 	} `json:"goal_params"`
	// 	ID      int64 `json:"id"`
	// 	IsAdmin bool  `json:"is_admin"`
	// 	Rank    int64 `json:"rank"`
	// 	Status  int64 `json:"status"`
	// 	User    struct {
	// 		HighlightedPhotoID int64  `json:"highlighted_photo_id"`
	// 		ID                 int64  `json:"id"`
	// 		Name               string `json:"name"`
	// 	} `json:"user"`
	// } `json:"relevant_goal_participants"`
	// SelfMembershipID int64 `json:"self_membership_id"`
	// SiteID           int64 `json:"site_id"`
	// SlimFavorites    []struct {
	// 	AssociatedObjectID   int64  `json:"associated_object_id"`
	// 	AssociatedObjectType string `json:"associated_object_type"`
	// 	ID                   int64  `json:"id"`
	// } `json:"slim_favorites"`
	TimeZone string `json:"time_zone"`
	// TotalRouteDistance         float64 `json:"total_route_distance"`
	// TripsIncludedInTotalsCount int64   `json:"trips_included_in_totals_count"`
	// UnseenUpdates              []struct {
	// 	Count int64  `json:"count"`
	// 	Key   string `json:"key"`
	// } `json:"unseen_updates"`
	// UpdatedAt    string               `json:"updated_at"`
	// UserSummary  map[string][]float64 `json:"user_summary"`
	// Visibility   int64                `json:"visibility"`
	// VO2Max       float64              `json:"vo2max"`
	// WeeksStartOn int64                `json:"weeks_start_on"`
	// Weight       float64              `json:"weight"`
}

// type Coordinate struct {
// 	Latitude  float64 `json:"lat"`
// 	Longitude float64 `json:"lng"`
// }

// type Photo struct {
// 	ID                int       `json:"id"`
// 	GroupMembershipID int       `json:"group_membership_id"`
// 	Caption           string    `json:"caption"`
// 	CreatedAt         string    `json:"created_at"`
// 	Position          int       `json:"position"`
// 	Visibility        int       `json:"visibility"`
// 	Latitude          float64   `json:"lat"`
// 	Longitude         float64   `json:"lng"`
// 	Published         bool      `json:"published"`
// 	CapturedAt        time.Time `json:"captured_at"`
// 	UserID            int       `json:"user_id"`
// 	UpdatedAt         string    `json:"updated_at"`
// 	Width             int       `json:"width"`
// 	Height            int       `json:"height"`
// 	OptionalUUID      string    `json:"optional_uuid"`
// 	Checksum          string    `json:"checksum"`
// }

// type MetricsSummary struct {
// 	Max int     `json:"max"`
// 	Min float64 `json:"min"`
// 	Avg float64 `json:"avg"`
// }

// type Metrics struct {
// 	ID             int             `json:"id"`
// 	ParentID       int             `json:"parent_id"`
// 	ParentType     string          `json:"parent_type"`
// 	CreatedAt      string          `json:"created_at"`
// 	UpdatedAt      string          `json:"updated_at"`
// 	Distance       float64         `json:"distance"`
// 	StartElevation float64         `json:"startElevation"`
// 	EndElevation   float64         `json:"endElevation"`
// 	NumPoints      int             `json:"numPoints"`
// 	EleGain        float64         `json:"ele_gain"`
// 	EleLoss        float64         `json:"ele_loss"`
// 	V              int             `json:"v"`
// 	Elevation      *MetricsSummary `json:"ele"`
// 	Grade          *MetricsSummary `json:"grade"`
// 	Watts          *MetricsSummary `json:"watts"`
// 	Cadence        *MetricsSummary `json:"cad"`
// 	HeartRate      *MetricsSummary `json:"hr"`
// }

// type Route struct {
// 	Type  string `json:"type"`
// 	Route struct {
// 		ID                       int           `json:"id"`
// 		HighlightedPhotoID       int           `json:"highlighted_photo_id"`
// 		HighlightedPhotoChecksum string        `json:"highlighted_photo_checksum"`
// 		Distance                 float64       `json:"distance"`
// 		ElevationGain            float64       `json:"elevation_gain"`
// 		ElevationLoss            float64       `json:"elevation_loss"`
// 		TrackID                  string        `json:"track_id"`
// 		UserID                   int           `json:"user_id"`
// 		PavementType             string        `json:"pavement_type"`
// 		PavementTypeID           int           `json:"pavement_type_id"`
// 		RecreationTypeIds        []interface{} `json:"recreation_type_ids"`
// 		Visibility               int           `json:"visibility"`
// 		CreatedAt                string        `json:"created_at"`
// 		UpdatedAt                string        `json:"updated_at"`
// 		Name                     string        `json:"name"`
// 		Description              string        `json:"description"`
// 		FirstLng                 float64       `json:"first_lng"`
// 		FirstLat                 float64       `json:"first_lat"`
// 		LastLat                  float64       `json:"last_lat"`
// 		LastLng                  float64       `json:"last_lng"`
// 		BoundingBox              *[]Coordinate `json:"bounding_box"`
// 		Locality                 string        `json:"locality"`
// 		PostalCode               string        `json:"postal_code"`
// 		AdministrativeArea       string        `json:"administrative_area"`
// 		CountryCode              string        `json:"country_code"`
// 		PrivacyCode              interface{}   `json:"privacy_code"`
// 		User                     struct {
// 			ID                       int     `json:"id"`
// 			CreatedAt                string  `json:"created_at"`
// 			Description              string  `json:"description"`
// 			Interests                string  `json:"interests"`
// 			Locality                 string  `json:"locality"`
// 			AdministrativeArea       string  `json:"administrative_area"`
// 			AccountLevel             int     `json:"account_level"`
// 			TotalTripDistance        float64 `json:"total_trip_distance"`
// 			TotalTripDuration        int     `json:"total_trip_duration"`
// 			TotalTripElevationGain   float64 `json:"total_trip_elevation_gain"`
// 			Name                     string  `json:"name"`
// 			HighlightedPhotoID       int     `json:"highlighted_photo_id"`
// 			HighlightedPhotoChecksum string  `json:"highlighted_photo_checksum"`
// 		} `json:"user"`
// 		HasCoursePoints  bool               `json:"has_course_points"`
// 		NavEnabled       bool               `json:"nav_enabled"`
// 		Rememberable     bool               `json:"rememberable"`
// 		Metrics          *Metrics           `json:"metrics"`
// 		Photos           *[]Photo           `json:"photos"`
// 		TrackPoints      []*Point           `json:"track_points"`
// 		CoursePoints     []*Point           `json:"course_points"`
// 		PointsOfInterest []*PointOfInterest `json:"points_of_interest"`
// 	} `json:"route"`
// }

// type Point struct {
// 	Latitude    float64 `json:"x"`
// 	Longitude   float64 `json:"y"`
// 	Distance    int     `json:"d"`
// 	Elevation   float64 `json:"e"`
// 	Time        int64   `json:"t"`
// 	Cadence     int     `json:"c"`
// 	HeartRate   int     `json:"h"`
// 	Power       int     `json:"p"`
// 	Speed       int     `json:"s"`
// 	Description string  `json:"description,omitempty"`
// }

// type PointOfInterest struct {
// 	ID         int     `json:"id"`
// 	Longitude  float64 `json:"lng"`
// 	Latitude   float64 `json:"lat"`
// 	URL        string  `json:"url"`
// 	MongoID    string  `json:"mongo_id"`
// 	ParentID   int     `json:"parent_id"`
// 	ParentType string  `json:"parent_type"`
// 	CreatedAt  string  `json:"created_at"`
// 	UpdatedAt  string  `json:"updated_at"`
// 	V          int     `json:"v"`
// 	T          int     `json:"t"`
// 	N          string  `json:"n"`
// 	D          string  `json:"d"`
// 	UID        int     `json:"uid"`
// 	// Pids       []interface{} `json:"pids"`
// }

// SegmentMatches  []struct {
// 	ID             int         `json:"id"`
// 	CreatedAt      string      `json:"created_at"`
// 	UpdatedAt      string      `json:"updated_at"`
// 	MongoID        string      `json:"mongo_id"`
// 	UserID         int         `json:"user_id"`
// 	SegmentID      int         `json:"segment_id"`
// 	ParentType     string      `json:"parent_type"`
// 	ParentID       int         `json:"parent_id"`
// 	FinalTime      interface{} `json:"final_time"`
// 	Visibility     int         `json:"visibility"`
// 	StartIndex     int         `json:"start_index"`
// 	EndIndex       int         `json:"end_index"`
// 	Duration       interface{} `json:"duration"`
// 	MovingTime     interface{} `json:"moving_time"`
// 	AscentTime     interface{} `json:"ascent_time"`
// 	PersonalRecord interface{} `json:"personal_record"`
// 	Vam            interface{} `json:"vam"`
// 	StartedAt      interface{} `json:"started_at"`
// 	Distance       float64     `json:"distance"`
// 	AvgSpeed       interface{} `json:"avg_speed"`
// 	Rank           interface{} `json:"rank"`
// 	Segment        struct {
// 		Title   string `json:"title"`
// 		Slug    string `json:"slug"`
// 		ToParam string `json:"to_param"`
// 	} `json:"segment"`
// 	Metrics struct {
// 		ID          int             `json:"id"`
// 		ParentID    int             `json:"parent_id"`
// 		ParentType  string          `json:"parent_type"`
// 		CreatedAt   string          `json:"created_at"`
// 		UpdatedAt   string          `json:"updated_at"`
// 		Elevation   *MetricsSummary `json:"ele"`
// 		HeartRate   *MetricsSummary `json:"hr"`
// 		Cadence     *MetricsSummary `json:"cad"`
// 		Speed       *MetricsSummary `json:"speed"`
// 		Grade       *MetricsSummary `json:"grade"`
// 		Watts       *MetricsSummary `json:"watts"`
// 		Stationary  bool            `json:"stationary"`
// 		Duration    int             `json:"duration"`
// 		FirstTime   int             `json:"firstTime"`
// 		MovingTime  int             `json:"movingTime"`
// 		StoppedTime int             `json:"stoppedTime"`
// 		Pace        float64         `json:"pace"`
// 		MovingPace  float64         `json:"movingPace"`
// 		AscentTime  int             `json:"ascentTime"`
// 		DescentTime int             `json:"descentTime"`
// 		Vam         float64         `json:"vam"`
// 		// TripSummary struct {
// 		// 	Num0 []float64 `json:"0"`
// 		// } `json:"tripSummary"`
// 		// HrZones struct {
// 		// 	Num1 int `json:"1"`
// 		// 	Num2 int `json:"2"`
// 		// 	Num3 int `json:"3"`
// 		// 	Num4 int `json:"4"`
// 		// 	Num5 int `json:"5"`
// 		// } `json:"hr_zones"`
// 		Distance       float64 `json:"distance"`
// 		StartElevation float64 `json:"startElevation"`
// 		EndElevation   float64 `json:"endElevation"`
// 		NumPoints      int     `json:"numPoints"`
// 		EleGain        float64 `json:"ele_gain"`
// 		EleLoss        float64 `json:"ele_loss"`
// 		IsClimb        bool    `json:"isClimb"`
// 		UciScore       float64 `json:"uciScore"`
// 		UciCategory    int     `json:"uciCategory"`
// 		FietsIndex     float64 `json:"fietsIndex"`
// 		V              int     `json:"v"`
// 	} `json:"metrics,omitempty"`
// 	Metrics struct {
// 		ID         int    `json:"id"`
// 		ParentID   int    `json:"parent_id"`
// 		ParentType string `json:"parent_type"`
// 		CreatedAt  string `json:"created_at"`
// 		UpdatedAt  string `json:"updated_at"`
// 		Ele        struct {
// 			Max  float64 `json:"max"`
// 			Min  float64 `json:"min"`
// 			Min  float64 `json:"_min"`
// 			Max  float64 `json:"_max"`
// 			MaxI int     `json:"max_i"`
// 			Avg  float64 `json:"_avg"`
// 			Avg  float64 `json:"avg"`
// 		} `json:"ele"`
// 		Grade struct {
// 			Max  float64 `json:"max"`
// 			Min  float64 `json:"min"`
// 			Min  float64 `json:"_min"`
// 			Max  float64 `json:"_max"`
// 			MaxI int     `json:"max_i"`
// 			MinI int     `json:"min_i"`
// 			Avg  float64 `json:"_avg"`
// 			Avg  float64 `json:"avg"`
// 		} `json:"grade"`
// 		Distance       float64 `json:"distance"`
// 		StartElevation float64 `json:"startElevation"`
// 		EndElevation   float64 `json:"endElevation"`
// 		NumPoints      int     `json:"numPoints"`
// 		EleGain        float64 `json:"ele_gain"`
// 		EleLoss        int     `json:"ele_loss"`
// 		V              int     `json:"v"`
// 		Watts          struct {
// 		} `json:"watts"`
// 		Cad struct {
// 		} `json:"cad"`
// 		Hr struct {
// 		} `json:"hr"`
// 	} `json:"metrics,omitempty"`
// 	Metrics struct {
// 		ID         int    `json:"id"`
// 		ParentID   int    `json:"parent_id"`
// 		ParentType string `json:"parent_type"`
// 		CreatedAt  string `json:"created_at"`
// 		UpdatedAt  string `json:"updated_at"`
// 		Ele        struct {
// 			Max  float64 `json:"max"`
// 			Min  int     `json:"min"`
// 			Min  int     `json:"_min"`
// 			Max  float64 `json:"_max"`
// 			MinI int     `json:"min_i"`
// 			Avg  float64 `json:"_avg"`
// 			Avg  float64 `json:"avg"`
// 		} `json:"ele"`
// 		Grade struct {
// 			Max  float64 `json:"max"`
// 			Min  float64 `json:"min"`
// 			Min  float64 `json:"_min"`
// 			Max  float64 `json:"_max"`
// 			MaxI int     `json:"max_i"`
// 			MinI int     `json:"min_i"`
// 			Avg  float64 `json:"_avg"`
// 			Avg  float64 `json:"avg"`
// 		} `json:"grade"`
// 		Distance       float64 `json:"distance"`
// 		StartElevation float64 `json:"startElevation"`
// 		EndElevation   int     `json:"endElevation"`
// 		NumPoints      int     `json:"numPoints"`
// 		EleGain        float64 `json:"ele_gain"`
// 		EleLoss        float64 `json:"ele_loss"`
// 		V              int     `json:"v"`
// 		Watts          struct {
// 		} `json:"watts"`
// 		Cad struct {
// 		} `json:"cad"`
// 		Hr struct {
// 		} `json:"hr"`
// 	} `json:"metrics,omitempty"`
// 	Metrics struct {
// 		ID         int    `json:"id"`
// 		ParentID   int    `json:"parent_id"`
// 		ParentType string `json:"parent_type"`
// 		CreatedAt  string `json:"created_at"`
// 		UpdatedAt  string `json:"updated_at"`
// 		Ele        struct {
// 			Max  float64 `json:"max"`
// 			Min  float64 `json:"min"`
// 			Min  float64 `json:"_min"`
// 			Max  float64 `json:"_max"`
// 			MaxI int     `json:"max_i"`
// 			Avg  float64 `json:"_avg"`
// 			Avg  float64 `json:"avg"`
// 		} `json:"ele"`
// 		Grade struct {
// 			Max  float64 `json:"max"`
// 			Min  float64 `json:"min"`
// 			Min  float64 `json:"_min"`
// 			Max  float64 `json:"_max"`
// 			MaxI int     `json:"max_i"`
// 			MinI int     `json:"min_i"`
// 			Avg  float64 `json:"_avg"`
// 			Avg  float64 `json:"avg"`
// 		} `json:"grade"`
// 		Distance       float64 `json:"distance"`
// 		StartElevation float64 `json:"startElevation"`
// 		EndElevation   float64 `json:"endElevation"`
// 		NumPoints      int     `json:"numPoints"`
// 		EleGain        float64 `json:"ele_gain"`
// 		EleLoss        float64 `json:"ele_loss"`
// 		V              int     `json:"v"`
// 		Watts          struct {
// 		} `json:"watts"`
// 		Cad struct {
// 		} `json:"cad"`
// 		Hr struct {
// 		} `json:"hr"`
// 	} `json:"metrics,omitempty"`
// 	Metrics struct {
// 		ID         int    `json:"id"`
// 		ParentID   int    `json:"parent_id"`
// 		ParentType string `json:"parent_type"`
// 		CreatedAt  string `json:"created_at"`
// 		UpdatedAt  string `json:"updated_at"`
// 		Ele        struct {
// 			Max  float64 `json:"max"`
// 			Min  float64 `json:"min"`
// 			Min  float64 `json:"_min"`
// 			Max  float64 `json:"_max"`
// 			MinI int     `json:"min_i"`
// 			MaxI int     `json:"max_i"`
// 			Avg  float64 `json:"_avg"`
// 			Avg  float64 `json:"avg"`
// 		} `json:"ele"`
// 		Grade struct {
// 			Max  float64 `json:"max"`
// 			Min  float64 `json:"min"`
// 			Min  float64 `json:"_min"`
// 			Max  float64 `json:"_max"`
// 			MaxI int     `json:"max_i"`
// 			MinI int     `json:"min_i"`
// 			Avg  float64 `json:"_avg"`
// 			Avg  float64 `json:"avg"`
// 		} `json:"grade"`
// 		Distance       float64 `json:"distance"`
// 		StartElevation float64 `json:"startElevation"`
// 		EndElevation   float64 `json:"endElevation"`
// 		NumPoints      int     `json:"numPoints"`
// 		EleGain        float64 `json:"ele_gain"`
// 		EleLoss        float64 `json:"ele_loss"`
// 		V              int     `json:"v"`
// 		Watts          struct {
// 		} `json:"watts"`
// 		Cad struct {
// 		} `json:"cad"`
// 		Hr struct {
// 		} `json:"hr"`
// 	} `json:"metrics,omitempty"`
// } `json:"segment_matches"`
