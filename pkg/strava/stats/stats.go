package stats

import (
	"math"
	"sort"

	"github.com/martinlindhe/unit"

	"github.com/bzimmer/gravl/pkg/common/stats"
	"github.com/bzimmer/gravl/pkg/strava"
)

func (s *Stats) KOMs(acts []*strava.Activity) []*strava.SegmentEffort {
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

func (s *Stats) HourRecord(acts []*strava.Activity) *strava.Activity {
	return strava.ReduceActivityPtr(func(act0, act1 *strava.Activity) *strava.Activity {
		if act0.AverageSpeed > act1.AverageSpeed {
			return act0
		}
		return act1
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		var dst, spd float64
		switch s.Units {
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

func (s *Stats) Distances(acts []*strava.Activity) []int {
	var vals []int
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		var val float64
		switch s.Units {
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

func (s *Stats) Eddington(acts []*strava.Activity) stats.EddingtonNumber {
	return stats.Eddington(s.Distances(acts))
}

func (s *Stats) BenfordsLaw(acts []*strava.Activity) stats.Benford {
	return stats.BenfordsLaw(s.Distances(acts))
}

func (s *Stats) ClimbingNumber(acts []*strava.Activity) int {
	var cnt int
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		cnt++
		return true
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		dst := act.Distance
		elv := act.TotalElevationGain
		switch s.Units {
		case UnitsImperial:
			dst = (unit.Length(dst) * unit.Meter).Miles()
			elv = (unit.Length(elv) * unit.Meter).Feet()
		case UnitsMetric:
			dst = (unit.Length(dst) * unit.Meter).Kilometers()
		}
		return int(elv/dst) > s.ClimbingNumberThreshold
	}, acts))
	return cnt
}

func (s *Stats) PythagoreanNumber(acts []*strava.Activity) []*PythagoreanNumber {
	var i int
	pn := make([]*PythagoreanNumber, len(acts))
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		// unit conversion is unnecessary as both distance and elevation are measured in meters
		pn[i] = &PythagoreanNumber{
			Activity: act,
			Number:   math.Sqrt(math.Pow(act.Distance, 2) + math.Pow(act.TotalElevationGain, 2)),
		}
		i++
		return true
	}, acts)
	sort.Slice(pn, func(i, j int) bool {
		return pn[i].Number > pn[j].Number
	})
	return pn
}
