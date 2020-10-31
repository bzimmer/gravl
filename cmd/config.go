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
	rwgpsUser      bool
	rwgpsAPIKey    string
	rwgpsAuthToken string

	// strava
	stravaAPIKey       string
	stravaAPISecret    string
	stravaAccessToken  string
	stravaRefreshToken string

	// visual crossing
	visualcrossingAPIKey string

	// root
	debug      bool
	compact    bool
	monochrome bool
	verbosity  string
	config     string
	encoder    *json.Encoder
)
