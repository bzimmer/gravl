package analysis

import (
	"time"

	"github.com/bzimmer/gravl/pkg/activity/strava"
)

type Activity struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	StartDate    time.Time `json:"startdate"`
	Distance     float64   `json:"distance"`
	Elevation    float64   `json:"elevation"`
	AverageSpeed float64   `json:"average_speed"`
	Type         string    `json:"type"`
}

func ToActivity(act *strava.Activity, units Units) *Activity {
	var dst, elv, spd float64
	switch units {
	case Metric:
		dst = act.Distance.Kilometers()
		elv = act.ElevationGain.Meters()
		spd = act.AverageSpeed.KilometersPerHour()
	case Imperial:
		dst = act.Distance.Miles()
		elv = act.ElevationGain.Feet()
		spd = act.AverageSpeed.MilesPerHour()
	}
	return &Activity{
		ID:           act.ID,
		Name:         act.Name,
		StartDate:    act.StartDate,
		Distance:     dst,
		Elevation:    elv,
		AverageSpeed: spd,
		Type:         act.Type,
	}
}
