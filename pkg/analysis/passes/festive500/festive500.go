package festive500

import (
	"math"
	"sort"
	"time"

	"github.com/martinlindhe/unit"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

const doc = `festive500 returns the activities and distance ridden during the annual #festive500 challenge

Thanks Rapha! https://www.rapha.cc/us/en_US/stories/festive-500`

var activityTypes = map[string]bool{
	"Ride":        true,
	"VirtualRide": true,
	"Handcycle":   true,
}

type Result struct {
	Activities        []*analysis.Activity `json:"activities"`
	DistanceCompleted float64              `json:"completed"`
	DistanceRemaining float64              `json:"remaining"`
	Percent           float64              `json:"percent"`
	Success           bool                 `json:"success"`
	MovingTime        unit.Duration        `json:"moving_time"`
}

func run(ctx *analysis.Context, pass []*strava.Activity) (interface{}, error) {
	var dst float64
	var dur unit.Duration
	var acts []*analysis.Activity
	for i := 0; i < len(pass); i++ {
		act := pass[i]
		_, ok := activityTypes[act.Type]
		if !ok {
			continue
		}
		_, month, date := act.StartDateLocal.Date()
		ok = (month == time.December && date >= 24 && date <= 31)
		if ok {
			dur = dur + act.MovingTime
			dst = dst + act.Distance.Kilometers()
			acts = append(acts, analysis.ToActivity(act, analysis.Metric))
		}
	}
	sort.Slice(acts, func(i, j int) bool {
		return acts[i].StartDate.Before(acts[j].StartDate)
	})
	return &Result{
		Activities:        acts,
		DistanceCompleted: dst,
		DistanceRemaining: math.Max(500.0-dst, 0),
		Percent:           (dst / 500.0) * 100,
		Success:           dst >= 500.0,
		MovingTime:        dur}, nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "festive500",
		Doc:  doc,
		Run:  run,
	}
}
