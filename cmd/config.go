package cmd

import "encoding/json"

var (
	// serve
	port       int
	origin     string
	sessionKey string

	// rwgps
	rwgpsTrip      bool
	rwgpsRoute     bool
	rwgpsAthlete   bool
	rwgpsAPIKey    string
	rwgpsAuthToken string

	// strava
	stravaRoute        bool
	stravaAthlete      bool
	stravaRefresh      bool
	stravaActivity     bool
	stravaAPIKey       string
	stravaAPISecret    string
	stravaAccessToken  string
	stravaRefreshToken string

	// visual crossing
	visualcrossingAPIKey string

	// root
	compact     bool
	monochrome  bool
	httptracing bool
	verbosity   string
	config      string

	// encoding
	encoder *json.Encoder
	decoder *json.Decoder
)
