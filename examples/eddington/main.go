package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/davecgh/go-spew/spew"
	"github.com/logic-building/functional-go/set"
	"github.com/martinlindhe/unit"
	"github.com/valyala/fastjson"

	"github.com/bzimmer/gravl/pkg/common/stats"
	"github.com/bzimmer/gravl/pkg/strava"
)

const (
	year              = 2020
	imperial          = true
	commutes          = false
	climbingThreshold = 0.0
)

var (
	// rides = set.NewStr([]string{"Ride", "VirtualRide"})
	rides = set.NewStr([]string{"Ride"})
)

func MustReadActivities(filename string) map[int][]*strava.Activity {
	var err error
	var sc fastjson.Scanner
	var acts []*strava.Activity

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	sc.InitBytes(b)
	for sc.Next() {
		if err = sc.Error(); err != nil {
			panic(err)
		}
		val := sc.Value()
		act := &strava.Activity{}
		err = json.Unmarshal(val.MarshalTo(nil), act)
		if err != nil {
			panic(err)
		}
		acts = append(acts, act)
	}

	return strava.GroupByIntActivityPtr(func(act *strava.Activity) int {
		return act.StartDateLocal.Year()
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		return rides.Contains(act.Type)
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		if year == 0 {
			return true
		}
		return year == act.StartDateLocal.Year()
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		if !commutes {
			return !act.Commute
		}
		return true
	}, acts))))
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
		var dst, spd float64
		if imperial {
			dst = (unit.Length(act.Distance) * unit.Meter).Miles()
			spd = (unit.Speed(act.AverageSpeed) * unit.KilometersPerHour).MilesPerHour()
		} else {
			dst = (unit.Length(act.Distance) * unit.Meter).Kilometers()
			spd = (unit.Speed(act.AverageSpeed) * unit.KilometersPerHour).KilometersPerHour()
		}
		return dst >= spd
	}, acts))
}

func Distances(acts []*strava.Activity) []int {
	var vals []int
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		var val float64
		if imperial {
			val = (unit.Length(act.Distance) * unit.Meter).Miles()
		} else {
			val = (unit.Length(act.Distance) * unit.Meter).Kilometers()
		}
		vals = append(vals, int(val))
		return true
	}, acts)
	return vals
}

func Eddington(acts []*strava.Activity) stats.EddingtonNumber {
	return stats.Eddington(Distances(acts))
}

func BenfordsLaw(acts []*strava.Activity) stats.Benford {
	return stats.BenfordsLaw(Distances(acts))
}

func ClimbingNumber(acts []*strava.Activity) int {
	var cnt int
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		cnt++
		return true
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		dst := act.Distance
		elv := act.TotalElevationGain
		if imperial {
			dst = (unit.Length(dst) * unit.Meter).Miles()
			elv = (unit.Length(elv) * unit.Meter).Feet()
		} else {
			dst = (unit.Length(dst) * unit.Meter).Kilometers()
		}
		return dst > climbingThreshold && elv/dst > 100
	}, acts))
	return cnt
}

func main() {
	var years []int
	acts := MustReadActivities(os.Args[1])
	for y := range acts {
		years = append(years, y)
	}
	sort.Ints(years)
	for _, year := range years {
		group := acts[year]
		fmt.Printf("\n%d (activities: %d)\n", year, len(group))
		fmt.Printf("Eddington: %d\n", Eddington(group).Number)
		act := HourRecord(group)
		fmt.Printf("HourRecord: %s (%d)\n", act.Name, act.ID)
		fmt.Printf("ClimbingNumber: %d\n", ClimbingNumber(group))
		koms := KOMs(group)
		sort.Slice(koms, func(i, j int) bool {
			return koms[i].Segment.ClimbCategory > koms[j].Segment.ClimbCategory
		})
		fmt.Printf("KOMs: %d\n", len(koms))
		for _, effort := range koms {
			seg := effort.Segment
			fmt.Printf("  > %s (avg %0.2f) (max %0.2f) (category %d)\n",
				seg.Name, seg.AverageGrade, seg.MaximumGrade, seg.ClimbCategory)
		}
		if len(group) > 25 {
			spewer := &spew.ConfigState{Indent: " ", SortKeys: true}
			spewer.Dump(BenfordsLaw(group))
			// fmt.Printf("Benford's Law: %0.4f\n", BenfordsLaw(group).ChiSquared)
		}
	}
}
