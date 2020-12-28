package festive500

import (
	"context"
	"sort"
	"time"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/strava/analysis"
)

const Doc = ``

var activityTypes = map[string]bool{
	"Ride":        true,
	"VirtualRide": true,
	"Handcycle":   true,
}

type Result struct {
	Activities []*analysis.Activity `json:"activities"`
	Distance   float64              `json:"distance"`
	Complete   float64              `json:"complete"`
	Success    bool                 `json:"success"`
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
	return &Result{
		Activities: acts,
		Distance:   dst,
		Complete:   (dst / 500.0) * 100,
		Success:    dst >= 500.0}, nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "festive500",
		Doc:  Doc,
		Run:  Run,
	}
}
