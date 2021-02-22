package ageride

import (
	"errors"
	"flag"
	"sort"
	"time"

	"github.com/rs/zerolog/log"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

const doc = `ageride returns all activities whose distance is greater than the athlete's age at the time of the activity`

const yearSeconds = 365.2425 /*(days)*/ * 24 /*(hours/day)*/ * 3600 /*(seconds/hour)*/

type ageRide struct {
	birthday *analysis.TimeFlag
}

type Result struct {
	Activities     []*analysis.Activity `json:"activities"`
	Count          int                  `json:"count"`
	DistanceMean   float64              `json:"distance_average"`
	DistanceMedian float64              `json:"distance_median"`
	DistanceTotal  float64              `json:"distance_total"`
}

func (a *ageRide) run(ctx *analysis.Context, pass []*strava.Activity) (interface{}, error) {
	var dsts []float64
	var acts []*analysis.Activity
	var bday = a.birthday.Get().(time.Time)
	if bday.IsZero() {
		return nil, errors.New("birthday not set")
	}
	log.Info().Time("birthday", bday).Msg("ageride")
	for i := 0; i < len(pass); i++ {
		act := pass[i]
		yrs := act.StartDateLocal.Sub(bday).Seconds() / yearSeconds
		switch ctx.Units {
		case analysis.Imperial:
			if act.Distance.Miles() > yrs {
				acts = append(acts, analysis.ToActivity(act, ctx.Units))
				dsts = append(dsts, act.Distance.Miles())
			}
		case analysis.Metric:
			if act.Distance.Kilometers() > yrs {
				acts = append(acts, analysis.ToActivity(act, ctx.Units))
				dsts = append(dsts, act.Distance.Kilometers())
			}
		}
	}
	var mean, median, total float64
	if len(dsts) > 0 {
		sort.Float64s(dsts)
		mean = stat.Mean(dsts, nil)
		median = stat.Quantile(0.5, stat.Empirical, dsts, nil)
		total = floats.Sum(dsts)
	}
	return &Result{
		Activities:     acts,
		Count:          len(acts),
		DistanceMean:   mean,
		DistanceMedian: median,
		DistanceTotal:  total}, nil
}

func New() *analysis.Analyzer {
	r := &ageRide{birthday: &analysis.TimeFlag{Time: time.Time{}}}
	fs := flag.NewFlagSet("ageride", flag.ExitOnError)
	fs.Var(r.birthday, "birthday", "the athlete's birthday in `YYYY-MM-DD` format")
	return &analysis.Analyzer{
		Name:  fs.Name(),
		Doc:   doc,
		Flags: fs,
		Run:   r.run,
	}
}
