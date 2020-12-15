package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"

	"github.com/davecgh/go-spew/spew"
	"github.com/logic-building/functional-go/set"
	"github.com/martinlindhe/unit"
	"github.com/valyala/fastjson"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/stats"
)

const (
	Year     = 2020
	Commutes = false
)

var (
	rides = set.NewStr([]string{
		// "Hike",
		"Ride",
		// "VirtualRide"
	})
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
		// if year == 0 {
		// 	return 0
		// }
		return act.StartDateLocal.Year()
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		return rides.Contains(act.Type)
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		// if year == 0 {
		// 	return true
		// }
		// return year == act.StartDateLocal.Year()
		return true
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		if !Commutes {
			return !act.Commute
		}
		return true
	}, acts))))
}

func main() {
	var years []int
	acts := MustReadActivities(os.Args[1])
	for y := range acts {
		years = append(years, y)
	}
	sort.Ints(years)

	// set up for imperial
	c := &stats.Config{
		Units:                   stats.UnitsImperial,
		ClimbingNumberThreshold: 100,
	}
	// c := stats.DefaultConfig
	for _, year := range years {
		group := acts[year]
		fmt.Printf("\n%d (activities: %d)\n", year, len(group))
		fmt.Printf("Eddington Number: %d\n", c.Eddington(group).Number)
		fmt.Printf("Climbing Number: %d\n", c.ClimbingNumber(group))
		act := c.HourRecord(group)
		fmt.Printf("Hour Record: %s (%d)\n", act.Name, act.ID)
		koms := c.KOMs(group)
		sort.Slice(koms, func(i, j int) bool {
			return koms[i].Segment.ClimbCategory > koms[j].Segment.ClimbCategory
		})
		fmt.Printf("KOMs: %d\n", len(koms))
		for _, effort := range koms {
			seg := effort.Segment
			fmt.Printf("  > %s (avg %0.2f) (max %0.2f) (category %d)\n",
				seg.Name, seg.AverageGrade, seg.MaximumGrade, seg.ClimbCategory)
		}
		pyt := c.Pythagorean(group)
		n := int(math.Min(10, float64(len(pyt))))
		fmt.Printf("Pythagorean (top %02d):\n", n)
		for i := 0; i < n; i++ {
			act = pyt[i]
			dst := (unit.Length(act.Distance) * unit.Meter).Miles()
			elv := (unit.Length(act.TotalElevationGain) * unit.Meter).Feet()
			fmt.Printf("  > %s (distance %0.2f), (elevation %0.2f)\n",
				act.Name, dst, elv)
		}
		if len(group) > 25 {
			spewer := &spew.ConfigState{Indent: " ", SortKeys: true}
			spewer.Dump(c.BenfordsLaw(group))
		}
	}
}
