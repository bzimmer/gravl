package analysis

import (
	"time"

	"github.com/bzimmer/gravl/pkg/strava"
)

type Activity struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"startdate"`
	Distance  float64   `json:"distance"`
	Elevation float64   `json:"elevation"`
	Type      string    `json:"type"`
}

func ToActivity(act *strava.Activity) *Activity {
	return ToActivityWithUnits(act, Imperial)
}

func ToActivityWithUnits(act *strava.Activity, units Units) *Activity {
	var dst, elv float64
	switch units {
	case Metric:
		dst = act.Distance.Kilometers()
		elv = act.ElevationGain.Meters()
	case Imperial:
		dst = act.Distance.Miles()
		elv = act.ElevationGain.Feet()
	}
	return &Activity{
		ID:        act.ID,
		Name:      act.Name,
		StartDate: act.StartDate,
		Distance:  dst,
		Elevation: elv,
	}
}
