package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/logic-building/functional-go/set"
	"github.com/martinlindhe/unit"
	"github.com/valyala/fastjson"

	"github.com/bzimmer/gravl/pkg/common/stats"
	"github.com/bzimmer/gravl/pkg/strava"
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
	return acts
}

func main() {
	var distances []int
	types := make((map[string]int))
	rides := set.NewStr([]string{"Ride", "VirtualRide"})

	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		// convert meters to miles
		miles := (unit.Length(act.Distance) * unit.Meter).Miles()
		distances = append(distances, int(miles))
		return true
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		types[act.Type] = types[act.Type] + 1
		return rides.Contains(act.Type)
	}, strava.FilterActivityPtr(func(act *strava.Activity) bool {
		return act.StartDateLocal.Year() == 2020
	}, MustReadActivities(os.Args[1]))))

	spewer := &spew.ConfigState{Indent: " ", SortKeys: true}
	spewer.Dump(stats.Eddington(distances))
	spewer.Dump(types)
}
