package bunt

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/tidwall/buntdb"

	"github.com/bzimmer/gravl/pkg/analysis/store"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
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
func Open(path string) (store.SourceSink, error) {
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
		return nil, err
	}
	var act *strava.Activity
	act, err = unmarshal(val)
	if err != nil {
		return nil, err
	}
	return act, nil
}

// Activities returns a slice of (potentially incomplete) Activity instances
func (b *bunt) Activities(ctx context.Context) ([]*strava.Activity, error) {
	var err error
	var acts []*strava.Activity
	err = b.db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend("", func(_, val string) bool {
			var act *strava.Activity
			if act, err = unmarshal(val); err != nil {
				return false
			}
			acts = append(acts, act)
			return true
		})
	})
	if err != nil {
		return nil, err
	}
	return acts, nil
}

// Save the activities to the source
func (b *bunt) Save(ctx context.Context, acts ...*strava.Activity) error {
	return b.db.Update(func(tx *buntdb.Tx) error {
		for i := range acts {
			act := acts[i]
			log.Debug().Int64("ID", act.ID).Str("name", act.Name).Msg("saving activity")
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
			log.Debug().Int64("id", acts[i].ID).Str("name", acts[i].Name).Msg("deleting activity")
			_, err := tx.Delete(id(acts[i].ID))
			if err != nil {
				return err
			}
		}
		return nil
	})
}
