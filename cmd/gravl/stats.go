package gravl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"sort"

	"github.com/davecgh/go-spew/spew"
	"github.com/martinlindhe/unit"
	"github.com/urfave/cli/v2"
	"github.com/valyala/fastjson"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/stats"
)

const (
	showBenfordsLaw = false
)

func readActivities(filename string) ([]*strava.Activity, error) {
	var err error
	var sc fastjson.Scanner
	var acts []*strava.Activity

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	sc.InitBytes(b)
	for sc.Next() {
		if err = sc.Error(); err != nil {
			return nil, err
		}
		val := sc.Value()
		act := &strava.Activity{}
		err = json.Unmarshal(val.MarshalTo(nil), act)
		if err != nil {
			return nil, err
		}
		acts = append(acts, act)
	}
	return acts, nil
}

func groupBy(c *cli.Context, acts []*strava.Activity) map[int][]*strava.Activity {
	// eliminate all commutes
	if !c.Bool("commutes") {
		acts = strava.FilterActivityPtr(func(act *strava.Activity) bool {
			return !act.Commute
		}, acts)
	}

	// filter activity types if specified
	activities := make(map[string]bool)
	activity := c.StringSlice("activity")
	for i := 0; i < len(activity); i++ {
		activities[activity[i]] = true
	}
	if len(activities) > 0 {
		acts = strava.FilterActivityPtr(func(act *strava.Activity) bool {
			return activities[act.Type]
		}, acts)
	}

	// filter years if specified
	years := make(map[int]bool)
	year := c.IntSlice("year")
	for i := 0; i < len(year); i++ {
		years[year[i]] = true
	}
	if len(years) > 0 {
		acts = strava.FilterActivityPtr(func(act *strava.Activity) bool {
			return years[act.StartDateLocal.Year()]
		}, acts)
	}

	// use a single grouping (possibly with fewer than all years) if totals are desired
	if c.Bool("totals") {
		return map[int][]*strava.Activity{0: acts}
	}

	// group away by year
	return strava.GroupByIntActivityPtr(func(act *strava.Activity) int {
		return act.StartDateLocal.Year()
	}, acts)
}

var statsCommand = &cli.Command{
	Name:     "stats",
	Category: "route",
	Usage:    "Compute stats from Strava activities",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "activity",
			Aliases: []string{"a"},
			Usage:   "Activity types to include",
		},
		&cli.IntSliceFlag{
			Name:    "year",
			Aliases: []string{"y"},
			Usage:   "Years to include, if not specified all years are calculated",
		},
		&cli.BoolFlag{
			Name:    "totals",
			Aliases: []string{"t"},
			Value:   false,
			Usage:   "Compute a total rather than grouped by years.",
		},
		&cli.BoolFlag{
			Name:    "commutes",
			Aliases: []string{"c"},
			Value:   false,
			Usage:   "Include commutes, filtered by default.",
		},
		&cli.BoolFlag{
			Name:    "metric",
			Aliases: []string{"m"},
			Value:   false,
			Usage:   "Use metric units (imperial is default).",
		},
	},
	Before: func(c *cli.Context) error {
		if c.NArg() == 0 {
			return errors.New("missing data file")
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		acts, err := readActivities(c.Args().First())
		if err != nil {
			return err
		}
		groups := groupBy(c, acts)
		var years []int
		for y := range groups {
			years = append(years, y)
		}
		sort.Ints(years)

		s := stats.ImperialStats
		if c.Bool("metric") {
			s = stats.MetricStats
		}
		for _, year := range years {
			group := groups[year]
			fmt.Printf("\n%d (activities: %d)\n", year, len(group))
			fmt.Printf("Eddington Number: %d\n", s.Eddington(group).Number)
			cn := s.ClimbingNumber(group)
			sort.Slice(cn, func(i, j int) bool {
				return cn[i].TotalElevationGain > cn[j].TotalElevationGain
			})
			fmt.Printf("Climbing Number: %d\n", len(cn))
			for i := 0; i < 10 && i < len(cn); i++ {
				act := cn[i]
				dst := act.Distance
				elv := act.TotalElevationGain
				if !c.Bool("metric") {
					dst = (unit.Length(dst) * unit.Meter).Miles()
					elv = (unit.Length(elv) * unit.Meter).Feet()
				} else {
					dst = (unit.Length(dst) * unit.Meter).Kilometers()
				}
				fmt.Printf("  > %s (distance %0.2f), (elevation %0.2f)\n",
					act.Name, dst, elv)
			}
			act := s.HourRecord(group)
			fmt.Printf("Hour Record: %s (%d)\n", act.Name, act.ID)
			koms := s.KOMs(group)
			sort.Slice(koms, func(i, j int) bool {
				return koms[i].Segment.ClimbCategory > koms[j].Segment.ClimbCategory
			})
			fmt.Printf("KOMs: %d\n", len(koms))
			for _, effort := range koms {
				seg := effort.Segment
				fmt.Printf("  > %s (%s) (avg %0.2f) (max %0.2f) (category %d)\n",
					seg.Name, seg.ActivityType, seg.AverageGrade, seg.MaximumGrade, seg.ClimbCategory)
			}
			pyt := s.PythagoreanNumber(group)
			n := int(math.Min(10, float64(len(pyt))))
			fmt.Printf("Pythagorean (top %02d):\n", n)
			for i := 0; i < n; i++ {
				pn := pyt[i]
				act := pn.Activity
				dst := act.Distance
				elv := act.TotalElevationGain
				if !c.Bool("metric") {
					dst = (unit.Length(dst) * unit.Meter).Miles()
					elv = (unit.Length(elv) * unit.Meter).Feet()
				} else {
					dst = (unit.Length(dst) * unit.Meter).Kilometers()
				}
				fmt.Printf("  > %s (number: %0.2f), (distance %0.2f), (elevation %0.2f)\n",
					act.Name, pn.Number, dst, elv)
			}
			if showBenfordsLaw && len(group) > 25 {
				spewer := &spew.ConfigState{Indent: " ", SortKeys: true}
				spewer.Dump(s.BenfordsLaw(group))
			}
		}
		return nil
	},
}
