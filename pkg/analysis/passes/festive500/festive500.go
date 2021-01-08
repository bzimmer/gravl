package festive500

import (
	"context"
	"sort"
	"time"

	"github.com/bzimmer/gravl/pkg/analysis"
)

const doc = ``

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

func run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var dst float64
	var acts []*analysis.Activity
	for i := 0; i < len(pass.Activities); i++ {
		act := pass.Activities[i]
		_, ok := activityTypes[act.Type]
		if !ok {
			continue
		}
		_, month, date := act.StartDateLocal.Date()
		ok = (month == time.December && date >= 24 && date <= 31)
		if ok {
			dst = dst + act.Distance.Kilometers()
			acts = append(acts, analysis.ToActivity(act, analysis.Metric))
		}
	}
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
		Doc:  doc,
		Run:  run,
	}
}
