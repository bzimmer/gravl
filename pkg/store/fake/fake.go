package fake

import (
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/rs/zerolog/log"

	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
	"github.com/bzimmer/gravl/pkg/store"
	"github.com/bzimmer/gravl/pkg/store/memory"
)

type fake struct {
	n    int
	fuzz bool
}

func (f *fake) Activities() ([]*strava.Activity, error) {
	activities := make([]*strava.Activity, f.n)
	defer func(t time.Time) {
		log.Info().
			Dur("elapsed", time.Since(t)).
			Bool("fuzz", f.fuzz).
			Int("activities", f.n).
			Msg("fake")
	}(time.Now())
	for i := 0; i < f.n; i++ {
		act, err := f.mk()
		if err != nil {
			return nil, err
		}
		act.ID = int64(i + 8200001)
		activities[i] = act
	}
	return activities, nil
}

func (f *fake) Close(acts map[int64]*strava.Activity) error {
	return nil
}

func Open(activities int, fuzz bool) (store.Store, error) {
	if fuzz {
		faker.SetIgnoreInterface(true)
		if err := faker.SetRandomMapAndSliceSize(2); err != nil {
			return nil, err
		}
	}
	return memory.Open(&fake{n: activities, fuzz: fuzz})
}

func (f *fake) mk() (*strava.Activity, error) {
	act := &strava.Activity{}
	if f.fuzz {
		if err := faker.FakeData(act); err != nil {
			return nil, err
		}
	}
	return act, nil
}
