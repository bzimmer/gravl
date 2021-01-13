package bolt

import (
	"context"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"

	"github.com/bzimmer/gravl/pkg/analysis/store"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

type bolt struct {
	db *bolthold.Store
}

// Open a bolt database; the file will be created if it does not exist
func Open(path string) (store.SourceSink, error) {
	directory := filepath.Dir(path)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		log.Info().Str("directory", directory).Msg("creating")
		if err = os.MkdirAll(directory, os.ModeDir|0700); err != nil {
			return nil, err
		}
	}
	db, err := bolthold.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}
	return &bolt{db: db}, nil
}

// Close the database
func (b *bolt) Close() error {
	return b.db.Close()
}

// Exists returns true if the activity exists, false otherwise
func (b *bolt) Exists(ctx context.Context, activityID int64) (bool, error) {
	_, err := b.Activity(ctx, activityID)
	if err != nil {
		if err == bolthold.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Activity returns a fully populated Activity
func (b *bolt) Activity(ctx context.Context, activityID int64) (*strava.Activity, error) {
	var act *strava.Activity
	if err := b.db.Get(activityID, &act); err != nil {
		return nil, err
	}
	return act, nil
}

// Activities returns a slice of (potentially incomplete) Activity instances
func (b *bolt) Activities(ctx context.Context) ([]*strava.Activity, error) {
	var acts []*strava.Activity
	err := b.db.ForEach(&bolthold.Query{}, func(act *strava.Activity) error {
		acts = append(acts, act)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return acts, nil
}

// Save the activities to the source
func (b *bolt) Save(ctx context.Context, acts ...*strava.Activity) error {
	return b.db.Bolt().Update(func(tx *bbolt.Tx) error {
		for i := range acts {
			act := acts[i]
			log.Debug().Int64("ID", act.ID).Str("name", act.Name).Msg("saving activity")
			if err := b.db.TxUpsert(tx, act.ID, act); err != nil {
				return err
			}
		}
		return nil
	})
}

// Remove the activities from the source
func (b *bolt) Remove(ctx context.Context, acts ...*strava.Activity) error {
	ids := make([]interface{}, len(acts))
	for i, act := range acts {
		ids[i] = act.ID
	}
	q := bolthold.Where("ID").In(ids...)
	log.Debug().Str("q", q.String()).Msg("deleting activities")
	return b.db.Bolt().Update(func(tx *bbolt.Tx) error {
		if err := b.db.TxDeleteMatching(tx, strava.Activity{}, q); err != nil {
			return err
		}
		return nil
	})
}
