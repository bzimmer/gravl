package stats

//go:generate stringer -type=Units -linecomment -output model_string.go

import (
	"math"
	"sort"

	"github.com/martinlindhe/unit"

	"github.com/bzimmer/gravl/pkg/common/stats"
	"github.com/bzimmer/gravl/pkg/strava"
)

type Units int

const (
	UnitsMetric   Units = iota // metric
	UnitsImperial              // imperial
)

type Config struct {
	Units                   Units
	ClimbingNumberThreshold int
}

// DefaultConfig configured for metric
var DefaultConfig = &Config{
	Units:                   UnitsMetric,
	ClimbingNumberThreshold: 20,
}

func (c *Config) KOMs(acts []*strava.Activity) []*strava.SegmentEffort {
	var efforts []*strava.SegmentEffort
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		for _, effort := range act.SegmentEfforts {
			for _, ach := range effort.Achievements {
				if ach.Rank == 1 && ach.Type == "overall" {
					efforts = append(efforts, effort)
					break
				}
			}
		}
		return true
	}, acts)
	return efforts
}

func KOMs(acts []*strava.Activity) []*strava.SegmentEffort {
	return DefaultConfig.KOMs(acts)
}

func (c *Config) HourRecord(acts []*strava.Activity) *strava.Activity {
	return strava.ReduceActivityPtr(func(act0, act1 *strava.Activity) *strava.Activity {
		if act0.AverageSpeed > act1.AverageSpeed {
			return act0
		}
		return act1
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		var dst, spd float64
		switch c.Units {
		case UnitsImperial:
			dst = (unit.Length(act.Distance) * unit.Meter).Miles()
			spd = (unit.Speed(act.AverageSpeed) * unit.KilometersPerHour).MilesPerHour()
		case UnitsMetric:
			// speed is already in kph
			dst = (unit.Length(act.Distance) * unit.Meter).Kilometers()
		}
		return dst >= spd
	}, acts))
}

func HourRecord(acts []*strava.Activity) *strava.Activity {
	return DefaultConfig.HourRecord(acts)
}

func (c *Config) Distances(acts []*strava.Activity) []int {
	var vals []int
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		var val float64
		switch c.Units {
		case UnitsImperial:
			val = (unit.Length(act.Distance) * unit.Meter).Miles()
		case UnitsMetric:
			val = (unit.Length(act.Distance) * unit.Meter).Kilometers()
		}
		vals = append(vals, int(val))
		return true
	}, acts)
	return vals
}

func Distances(acts []*strava.Activity) []int {
	return DefaultConfig.Distances(acts)
}

func (c *Config) Eddington(acts []*strava.Activity) stats.EddingtonNumber {
	return stats.Eddington(c.Distances(acts))
}

func Eddington(acts []*strava.Activity) stats.EddingtonNumber {
	return stats.Eddington(DefaultConfig.Distances(acts))
}

func (c *Config) BenfordsLaw(acts []*strava.Activity) stats.Benford {
	return stats.BenfordsLaw(c.Distances(acts))
}

func BenfordsLaw(acts []*strava.Activity) stats.Benford {
	return stats.BenfordsLaw(DefaultConfig.Distances(acts))
}

func (c *Config) ClimbingNumber(acts []*strava.Activity) int {
	var cnt int
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		cnt++
		return true
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		dst := act.Distance
		elv := act.TotalElevationGain
		switch c.Units {
		case UnitsImperial:
			dst = (unit.Length(dst) * unit.Meter).Miles()
			elv = (unit.Length(elv) * unit.Meter).Feet()
		case UnitsMetric:
			dst = (unit.Length(dst) * unit.Meter).Kilometers()
		}
		return int(elv/dst) > c.ClimbingNumberThreshold
	}, acts))
	return cnt
}

func ClimbingNumber(acts []*strava.Activity) int {
	return DefaultConfig.ClimbingNumber(acts)
}

func (c *Config) Pythagorean(acts []*strava.Activity) []*strava.Activity {
	m := make([]*strava.Activity, len(acts))
	copy(m, acts)
	sort.Slice(m, func(i, j int) bool {
		// unit conversion is unnecessary as both distance and elevation are measured in meters
		return (math.Pow(m[i].Distance, 2)+math.Pow(m[i].TotalElevationGain, 2) >
			math.Pow(m[j].Distance, 2)+math.Pow(m[j].TotalElevationGain, 2))
	})
	return m
}

func Pythagorean(acts []*strava.Activity) []*strava.Activity {
	return DefaultConfig.Pythagorean(acts)
}
