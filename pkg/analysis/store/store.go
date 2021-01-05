package store

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"

	"github.com/bzimmer/gravl/pkg/activity/strava"
)

type Store struct {
	store *bolthold.Store
}

func NewStore(store *bolthold.Store) *Store {
	return &Store{store: store}
}

func (s *Store) Update(ctx context.Context, source Source) (int, error) {
	var n int
	var err error
	log.Info().Msg("querying activities")
	acts, err := source.Activities(ctx)
	if err != nil {
		return n, err
	}
	log.Info().Int("n", len(acts)).Msg("found activities")
	err = s.store.Bolt().Update(func(tx *bbolt.Tx) error {
		for i := range acts {
			var t strava.Activity
			err = s.store.TxGet(tx, acts[i].ID, &t)
			if err == nil {
				continue
			}
			if err != nil && err != bolthold.ErrNotFound {
				return err
			}
			log.Info().Int64("ID", acts[i].ID).Msg("querying activity details")
			var act *strava.Activity
			act, err = source.Activity(ctx, acts[i].ID)
			if err != nil {
				return err
			}
			log.Info().Int64("ID", act.ID).Str("name", act.Name).Msg("saving activity details")
			if err = s.store.TxUpsert(tx, act.ID, act); err != nil {
				return err
			}
			n++
		}
		return nil
	})
	if err != nil {
		return n, err
	}
	return n, nil
}
