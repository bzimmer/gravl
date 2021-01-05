package festive500

import (
	"context"
	"sort"
	"time"

	"github.com/bzimmer/gravl/pkg/activity/strava"
	"github.com/bzimmer/gravl/pkg/analysis"
)

const Doc = ``

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
}

func Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var dst float64
	var acts []*analysis.Activity
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		_, ok := activityTypes[act.Type]
		if !ok {
			return true
		}
		_, month, date := act.StartDateLocal.Date()
		ok = (month == time.December && date >= 24 && date <= 31)
		if ok {
			dst = dst + act.Distance.Kilometers()
			acts = append(acts, analysis.ToActivity(act, analysis.Metric))
		}
		return true
	}, pass.Activities)
	sort.Slice(acts, func(i, j int) bool {
		return acts[i].StartDate.Before(acts[j].StartDate)
	})
	remaining := 500.0 - dst
	if remaining < 0 {
		remaining = 0
	}
	return &Result{
		Activities:        acts,
		DistanceCompleted: dst,
		DistanceRemaining: remaining,
		Percent:           (dst / 500.0) * 100,
		Success:           dst >= 500.0}, nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "festive500",
		Doc:  Doc,
		Run:  Run,
	}
}
