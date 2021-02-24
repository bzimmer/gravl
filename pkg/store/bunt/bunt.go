package bunt

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/tidwall/buntdb"

	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
	"github.com/bzimmer/gravl/pkg/store"
)

func id(activityID int64) string {
	return fmt.Sprintf("activity:%d", activityID)
}

func marshal(act *strava.Activity) (string, error) {
	val, err := json.Marshal(act)
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func unmarshal(val string) (*strava.Activity, error) {
	var act *strava.Activity
	if err := json.Unmarshal([]byte(val), &act); err != nil {
		return nil, err
	}
	return act, nil
}

type bunt struct {
	db *buntdb.DB
}

// Open a bolt database; the file will be created if it does not exist
func Open(path string) (store.Store, error) {
	log.Info().Str("path", path).Msg("bunt db")
	directory := filepath.Dir(path)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		log.Info().Str("directory", directory).Msg("creating")
		if err = os.MkdirAll(directory, os.ModeDir|0700); err != nil {
			return nil, err
		}
	}
	db, err := buntdb.Open(path)
	if err != nil {
		return nil, err
	}
	return &bunt{db: db}, nil
}

// Close the database
func (b *bunt) Close() error {
	return b.db.Close()
}

// Exists returns true if the activity exists, false otherwise
func (b *bunt) Exists(ctx context.Context, activityID int64) (bool, error) {
	err := b.db.View(func(tx *buntdb.Tx) error {
		_, err := tx.Get(id(activityID))
		return err
	})
	if err != nil {
		if err == buntdb.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Activity returns a fully populated Activity
func (b *bunt) Activity(ctx context.Context, activityID int64) (*strava.Activity, error) {
	var err error
	var val string
	err = b.db.View(func(tx *buntdb.Tx) error {
		val, err = tx.Get(id(activityID))
		return err
	})
	if err != nil {
		if err == buntdb.ErrNotFound {
			return nil, store.ErrNotFound
		}
		return nil, err
	}
	var act *strava.Activity
	act, err = unmarshal(val)
	if err != nil {
		return nil, err
	}
	return act, nil
}

// Activities returns a channel of activities and errors for an athlete
func (b *bunt) Activities(ctx context.Context) <-chan *strava.ActivityResult {
	acts := make(chan *strava.ActivityResult)
	go func() {
		defer close(acts)
		err := b.db.View(func(tx *buntdb.Tx) error {
			return tx.Ascend("", func(_, val string) bool {
				var r *strava.ActivityResult
				act, err := unmarshal(val)
				if err != nil {
					r = &strava.ActivityResult{Err: err}
				} else {
					r = &strava.ActivityResult{Activity: act}
				}
				select {
				case <-ctx.Done():
					log.Debug().Err(ctx.Err()).Msg("ctx is done")
					return false
				case acts <- r:
					return r.Err == nil
				}
			})
		})
		if err != nil {
			select {
			case <-ctx.Done():
				log.Debug().Err(ctx.Err()).Msg("ctx is done")
				return
			case acts <- &strava.ActivityResult{Err: err}:
				return
			}
		}
	}()
	return acts
}

// Save the activities to the source
func (b *bunt) Save(ctx context.Context, acts ...*strava.Activity) error {
	return b.db.Update(func(tx *buntdb.Tx) error {
		for i := range acts {
			act := acts[i]
			val, err := marshal(act)
			if err != nil {
				return err
			}
			if _, _, err := tx.Set(id(act.ID), val, nil); err != nil {
				return err
			}
		}
		return nil
	})
}

// Remove the activities from the source
func (b *bunt) Remove(ctx context.Context, acts ...*strava.Activity) error {
	return b.db.Update(func(tx *buntdb.Tx) error {
		for i := range acts {
			_, err := tx.Delete(id(acts[i].ID))
			if err != nil {
				return err
			}
		}
		return nil
	})
}
