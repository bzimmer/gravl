package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/logic-building/functional-go/set"
	"github.com/martinlindhe/unit"
	"github.com/valyala/fastjson"

	"github.com/bzimmer/gravl/pkg/common/stats"
	"github.com/bzimmer/gravl/pkg/strava"
)

const imperial = true
const commutes = false
const climbingThreshold = 0.0

var (
	year = 0
	// rides = set.NewStr([]string{"Ride", "VirtualRide"})
	rides = set.NewStr([]string{"Ride"})
	// spewer = &spew.ConfigState{Indent: " ", SortKeys: true}
)

func MustReadActivities(filename string) []*strava.Activity {
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

	return strava.FilterActivityPtr(func(act *strava.Activity) bool {
		return rides.Contains(act.Type)
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		if !commutes {
			return !act.Commute
		}
		return true
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		return year == 0 || act.StartDateLocal.Year() == year
	}, acts)))
}

func HourRecord(acts []*strava.Activity) *strava.Activity {
	return strava.ReduceActivityPtr(func(act0, act1 *strava.Activity) *strava.Activity {
		if act0.AverageSpeed > act1.AverageSpeed {
			return act0
		}
		return act1
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		miles := (unit.Length(act.Distance) * unit.Meter).Miles()
		mph := (unit.Speed(act.AverageSpeed) * unit.KilometersPerHour).MilesPerHour()
		return miles >= mph
	}, acts))
}

func Eddington(acts []*strava.Activity) stats.EddingtonNumber {
	var vals []int
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		val := act.Distance
		if imperial {
			val = (unit.Length(val) * unit.Meter).Miles()
		}
		vals = append(vals, int(val))
		return true
	}, acts)
	return stats.Eddington(vals)
}

func BenfordsLaw(acts []*strava.Activity) stats.Benford {
	var vals []int
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		val := act.Distance
		if imperial {
			val = (unit.Length(val) * unit.Meter).Miles()
		}
		vals = append(vals, int(val))
		return true
	}, acts)
	return stats.BenfordsLaw(vals)
}

func ClimbingNumber(acts []*strava.Activity) int {
	var cnt int
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		cnt++
		return true
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		miles := (unit.Length(act.Distance) * unit.Meter).Miles()
		elv := (unit.Length(act.TotalElevationGain) * unit.Meter).Feet()
		return miles > climbingThreshold && elv/miles > 100
	}, acts))
	return cnt
}

func GroupByIntActivityPtr(f func(act *strava.Activity) int, acts []*strava.Activity) map[int][]*strava.Activity {
	res := make(map[int][]*strava.Activity)
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		key := f(act)
		res[key] = append(res[key], act)
		return true
	}, acts)
	return res
}

func main() {
	var years []int
	acts := MustReadActivities(os.Args[1])
	groups := GroupByIntActivityPtr(func(act *strava.Activity) int {
		// return 0
		return act.StartDateLocal.Year()
	}, acts)
	for year := range groups {
		years = append(years, year)
	}
	sort.Ints(years)
	for _, year := range years {
		group := groups[year]
		fmt.Printf("\n%d (activities: %d)\n", year, len(group))
		fmt.Printf("Eddington: %d\n", Eddington(group).Number)
		act := HourRecord(group)
		fmt.Printf("HourRecord: %s (%d)\n", act.Name, act.ID)
		fmt.Printf("ClimbingNumber: %d\n", ClimbingNumber(group))
		if len(group) > 25 {
			fmt.Printf("Benford's Law: %0.4f\n", BenfordsLaw(group).ChiSquared)
		}
	}
}
