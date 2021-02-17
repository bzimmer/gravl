package koms

import (
	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

const doc = `koms returns all KOMs for the activities`

func run(ctx *analysis.Context, pass []*strava.Activity) (interface{}, error) {
	var efforts []*strava.SegmentEffort
	for i := 0; i < len(pass); i++ {
		act := pass[i]
		for _, effort := range act.SegmentEfforts {
			for _, ach := range effort.Achievements {
				if ach.Rank == 1 && ach.Type == "overall" {
					efforts = append(efforts, effort)
				}
			}
		}
	}
	return efforts, nil
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "koms",
		Doc:  doc,
		Run:  run,
	}
}
