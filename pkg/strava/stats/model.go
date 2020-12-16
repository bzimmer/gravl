package stats

import "github.com/bzimmer/gravl/pkg/strava"

//go:generate stringer -type=Units -linecomment -output model_string.go

type Units int

const (
	UnitsMetric   Units = iota // metric
	UnitsImperial              // imperial
)

type Stats struct {
	Units                   Units
	ClimbingNumberThreshold int
}

// MetricStats configured for metric units
var MetricStats = &Stats{
	Units:                   UnitsMetric,
	ClimbingNumberThreshold: 20,
}

// ImperialStats configured for imperial units
var ImperialStats = &Stats{
	Units:                   UnitsImperial,
	ClimbingNumberThreshold: 100,
}

// PythagoreanNumber for an activity
type PythagoreanNumber struct {
	Activity *strava.Activity `json:"activity"`
	Number   float64          `json:"number"`
}
