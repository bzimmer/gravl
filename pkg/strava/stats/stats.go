package stats

import (
	"math"
	"sort"

	"github.com/bzimmer/gravl/pkg/common/stats"
	"github.com/bzimmer/gravl/pkg/strava"
)

type Units int

const (
	Metric Units = iota
	Imperial
)

// Pythagorean number for an activity
type Pythagorean struct {
	Activity *strava.Activity `json:"activity"`
	Number   int              `json:"number"`
}

// Climbing number for an activity
type Climbing struct {
	Activity *strava.Activity `json:"activity"`
	Number   int              `json:"number"`
}

func KOMs(acts []*strava.Activity) []*strava.SegmentEffort {
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

func HourRecord(acts []*strava.Activity) *strava.Activity {
	return strava.ReduceActivityPtr(func(act0, act1 *strava.Activity) *strava.Activity {
		if act0.AverageSpeed > act1.AverageSpeed {
			return act0
		}
		return act1
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		dst := act.Distance
		spd := act.AverageSpeed
		return float64(dst) >= float64(spd)
	}, acts))
}

func Distances(acts []*strava.Activity, units Units) []int {
	var vals []int
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		var dst float64
		switch units {
		case Metric:
			dst = act.Distance.Kilometers()
		case Imperial:
			dst = act.Distance.Miles()
		}
		vals = append(vals, int(dst))
		return true
	}, acts)
	return vals
}

func EddingtonNumber(acts []*strava.Activity, units Units) stats.Eddington {
	return stats.EddingtonNumber(Distances(acts, units))
}

func BenfordsLaw(acts []*strava.Activity, units Units) stats.Benford {
	return stats.BenfordsLaw(Distances(acts, units))
}

func ClimbingNumber(acts []*strava.Activity, units Units, climbingThreshold int) []*Climbing {
	var climbings []*Climbing
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		var dst, elv float64
		switch units {
		case Metric:
			dst = act.Distance.Kilometers()
			elv = act.ElevationGain.Meters()
		case Imperial:
			dst = act.Distance.Miles()
			elv = act.ElevationGain.Feet()
		}
		c := int(elv / dst)
		if c > climbingThreshold {
			climbings = append(climbings, &Climbing{Activity: act, Number: c})
		}
		return true
	}, acts)
	return climbings
}

func PythagoreanNumber(acts []*strava.Activity) []*Pythagorean {
	var i int
	pn := make([]*Pythagorean, len(acts))
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		dst := act.Distance.Meters()
		elv := act.ElevationGain.Meters()
		n := int(math.Sqrt(math.Pow(dst, 2) + math.Pow(elv, 2)))
		pn[i] = &Pythagorean{Activity: act, Number: n}
		i++
		return true
	}, acts)
	sort.Slice(pn, func(i, j int) bool {
		return pn[i].Number > pn[j].Number
	})
	return pn
}
