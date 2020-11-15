package wta

import "time"

// Region .
type Region struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Subregions *[]Region `json:"subregions,omitempty"`
}

// TripReport .
type TripReport struct {
	Reporter string    `json:"reporter"`
	Title    string    `json:"title"`
	Report   string    `json:"report_url"`
	HikeDate time.Time `json:"hike_date"`
	Votes    int       `json:"votes"`
	Region   string    `json:"region"`
	Photos   int       `json:"photos"`
}

// TripReports .
type TripReports struct {
	Reporter string        `json:"reporter"`
	Reports  []*TripReport `json:"reports"`
}
