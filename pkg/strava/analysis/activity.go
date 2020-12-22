package analysis

//go:generate stringer -type=Units -linecomment -output=activity_string.go

import (
	"fmt"
	"time"

	"github.com/bzimmer/gravl/pkg/strava"
)

type Units int

const (
	Metric   Units = iota // metric
	Imperial              // imperial
)

type UnitsFlag struct {
	Units *Units
}

func (u *UnitsFlag) String() string {
	if u.Units == nil {
		// default to imperial
		return Imperial.String()
	}
	return u.Units.String()
}

func (u *UnitsFlag) Set(value string) error {
	switch value {
	case "imperial":
		*u.Units = Imperial
	case "metric":
		*u.Units = Metric
	default:
		return fmt.Errorf("unexpected unit '%s'", value)
	}
	return nil
}

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
