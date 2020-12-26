package store

import (
	"context"

	"github.com/rs/zerolog/log"
	bh "github.com/timshannon/bolthold"
	bolt "go.etcd.io/bbolt"

	"github.com/bzimmer/gravl/pkg/strava"
)

type Store struct {
	store *bh.Store
}

func NewStore(store *bh.Store) *Store {
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
	err = s.store.Bolt().Update(func(tx *bolt.Tx) error {
		var act *strava.Activity
		for i := range acts {
			err = s.store.TxGet(tx, acts[i].ID, act)
			if err != bh.ErrNotFound {
				continue
			}
			log.Info().Int64("ID", acts[i].ID).Msg("querying activity details")
			act, err = source.Activity(ctx, acts[i].ID)
			if err != nil {
				return err
			}
			log.Info().Int64("ID", act.ID).Msg("saving activity details")
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