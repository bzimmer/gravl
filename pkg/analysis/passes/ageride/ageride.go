package ageride

import (
	"context"
	"errors"
	"flag"
	"sort"
	"time"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/rs/zerolog/log"
)

const doc = ``

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

func (a *ageRide) run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var dsts []float64
	var acts []*analysis.Activity
	var bday = a.birthday.Get().(time.Time)
	if bday.IsZero() {
		return nil, errors.New("birthday not set")
	}
	log.Info().Time("birthday", bday).Msg("ageride")
	for i := 0; i < len(pass.Activities); i++ {
		act := pass.Activities[i]
		yrs := act.StartDateLocal.Sub(bday).Seconds() / yearSeconds
		switch pass.Units {
		case analysis.Imperial:
			if act.Distance.Miles() > yrs {
				acts = append(acts, analysis.ToActivity(act, pass.Units))
				dsts = append(dsts, act.Distance.Miles())
			}
		case analysis.Metric:
			if act.Distance.Kilometers() > yrs {
				acts = append(acts, analysis.ToActivity(act, pass.Units))
				dsts = append(dsts, act.Distance.Kilometers())
			}
		}
	}
	sort.Float64s(dsts)
	return &Result{
		Activities:     acts,
		Count:          len(acts),
		DistanceMean:   stat.Mean(dsts, nil),
		DistanceMedian: stat.Quantile(0.5, stat.Empirical, dsts, nil),
		DistanceTotal:  floats.Sum(dsts)}, nil
}

func New() *analysis.Analyzer {
	r := &ageRide{birthday: &analysis.TimeFlag{Time: time.Time{}}}
	fs := flag.NewFlagSet("ageride", flag.ExitOnError)
	// @todo(bzimmer) ideally this would use the athlete's birthdate on strava
	fs.Var(r.birthday, "birthday", "the athlete's birthday")
	return &analysis.Analyzer{
		Name:  fs.Name(),
		Doc:   doc,
		Flags: fs,
		Run:   r.run,
	}
}
