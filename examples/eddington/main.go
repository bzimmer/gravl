package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/logic-building/functional-go/set"
	"github.com/martinlindhe/unit"
	"github.com/valyala/fastjson"

	"github.com/bzimmer/gravl/pkg/common/stats"
	"github.com/bzimmer/gravl/pkg/strava"
)

const year = 2017
const imperial = true

var (
	rides = set.NewStr([]string{"Ride", "VirtualRide"})
	// rides  = set.NewStr([]string{"Ride"})
	spewer = &spew.ConfigState{Indent: " ", SortKeys: true}
)

type EddingtonFunc func(*strava.Activity) (int, bool)

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

	return acts
}

func Distance(act *strava.Activity) (int, bool) {
	// convert meters to miles
	val := act.Distance
	if imperial {
		val = (unit.Length(val) * unit.Meter).Miles()
	}
	return int(val), true
}

func Elevation(act *strava.Activity) (int, bool) {
	// convert meters to feet
	val := act.TotalElevationGain
	if imperial {
		val = (unit.Length(val) * unit.Meter).Feet()
	}
	return int(val), true
}

func Eddington(f EddingtonFunc, acts []*strava.Activity) stats.EddingtonNumber {
	var vals []int
	types := make((map[string]int))
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		val, ok := f(act)
		if ok {
			vals = append(vals, val)
		}
		return ok
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		types[act.Type] = types[act.Type] + 1
		return rides.Contains(act.Type)
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		return act.StartDateLocal.Year() == year
	}, acts)))
	return stats.Eddington(vals)
}

func BenfordsLaw(acts []*strava.Activity) stats.Benford {
	var vals []int
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		miles := (unit.Length(act.Distance) * unit.Meter).Miles()
		vals = append(vals, int(miles))
		return true
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		return rides.Contains(act.Type)
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		return act.StartDateLocal.Year() == year
	}, acts)))
	return stats.BenfordsLaw(vals)
}

func main() {
	acts := MustReadActivities(os.Args[1])
	for i, f := range []EddingtonFunc{Distance, Elevation} {
		ed := Eddington(f, acts)
		fmt.Printf("%d > %d\n", i, ed.Number)
	}
	spewer.Dump(BenfordsLaw(acts))
}
