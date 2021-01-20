package analysis

import (
	"time"

	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

type Activity struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	StartDate    time.Time `json:"startdate"`
	Distance     float64   `json:"distance"`
	Elevation    float64   `json:"elevation"`
	MovingTime   float64   `json:"movingtime"`
	AverageSpeed float64   `json:"averagespeed"`
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
		MovingTime:   act.MovingTime.Seconds(),
		Type:         act.Type,
	}
}
