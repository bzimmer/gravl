package stats

import (
	"time"

	"github.com/bzimmer/gravl/pkg/common/stats"
	"github.com/bzimmer/gravl/pkg/strava"
)

type Activity struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"startdate"`
	Distance  float64   `json:"distance"`
	Elevation float64   `json:"elevation"`
	Type      string    `json:"type"`
}

type PythagoreanAnalysis struct {
	Activity *Activity `json:"activity"`
	Number   int       `json:"number"`
}

type ClimbingAnalysis struct {
	Activity *Activity `json:"activity"`
	Number   int       `json:"number"`
}

type Analysis struct {
	Pythagorean []*PythagoreanAnalysis  `json:"pythagorean"`
	Climbing    []*ClimbingAnalysis     `json:"climbing"`
	HourRecord  *Activity               `json:"hour"`
	Eddington   *stats.Eddington        `json:"eddington"`
	Benford     *stats.Benford          `json:"benford"`
	KOMs        []*strava.SegmentEffort `json:"koms"`
}

type Analyzer struct {
	Activities        []*strava.Activity
	Units             Units
	ClimbingThreshold int
}

func (a *Analyzer) Analyze() *Analysis {
	dsts := Distances(a.Activities, a.Units)
	bd := stats.BenfordsLaw(dsts)
	ed := stats.EddingtonNumber(dsts)
	hr := HourRecord(a.Activities)

	var a2a = func(act *strava.Activity) *Activity {
		var dst, elv float64
		switch a.Units {
		case Metric:
			dst = act.Distance.Kilometers()
			elv = act.ElevationGain.Meters()
		case Imperial:
			dst = act.Distance.Miles()
			elv = act.ElevationGain.Feet()
		}
		return &Activity{
			ID:        act.ID,
			Name:      act.Name,
			StartDate: act.StartDate,
			Distance:  dst,
			Elevation: elv,
		}
	}

	pn := PythagoreanNumber(a.Activities)
	pna := make([]*PythagoreanAnalysis, len(pn))
	for i := 0; i < len(pn); i++ {
		pna[i] = &PythagoreanAnalysis{
			Activity: a2a(pn[i].Activity),
			Number:   pn[i].Number,
		}
	}

	cn := ClimbingNumber(a.Activities, a.Units, a.ClimbingThreshold)
	cna := make([]*ClimbingAnalysis, len(cn))
	for i := 0; i < len(cn); i++ {
		cna[i] = &ClimbingAnalysis{
			Activity: a2a(cn[i].Activity),
			Number:   cn[i].Number,
		}
	}

	ay := &Analysis{
		HourRecord:  a2a(hr),
		Eddington:   &ed,
		Benford:     &bd,
		Climbing:    cna,
		Pythagorean: pna,
	}
	return ay
}
