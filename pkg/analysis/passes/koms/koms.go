package koms

import (
	"context"

	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/strava"
)

const Doc = ``

func Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	var efforts []*strava.SegmentEffort
	strava.EveryActivityPtr(func(act *strava.Activity) bool {
		for _, effort := range act.SegmentEfforts {
			for _, ach := range effort.Achievements {
				if ach.Rank == 1 && ach.Type == "overall" {
					efforts = append(efforts, effort)
					break
				}
			}
		}
		return true
	}, pass.Activities)
	return efforts, nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "koms",
		Doc:  Doc,
		Run:  Run,
	}
}
