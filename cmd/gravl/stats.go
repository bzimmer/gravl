package gravl

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sort"

	"github.com/urfave/cli/v2"
	"github.com/valyala/fastjson"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/stats"
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
	// filter commutes if requested
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
	year := c.IntSlice("year")
	years := make(map[int]bool)
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
			Usage:   "Include commutes, (default: filtered).",
		},
		&cli.BoolFlag{
			Name:    "metric",
			Aliases: []string{"m"},
			Value:   false,
			Usage:   "Use metric units (default: imperial).",
		},
	},
	Before: func(c *cli.Context) error {
		if c.NArg() == 0 {
			return errors.New("missing data file")
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		units := stats.Imperial
		if c.Bool("metric") {
			units = stats.Metric
		}

		var threshold int
		switch units {
		case stats.Metric:
			threshold = 20
		case stats.Imperial:
			threshold = 100
		}

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

		res := make(map[int]*stats.Analysis)
		for _, year := range years {
			group := groups[year]
			anz := stats.Analyzer{
				Activities:        group,
				Units:             units,
				ClimbingThreshold: threshold,
			}
			res[year] = anz.Analyze()
		}
		return encoder.Encode(res)
	},
}
