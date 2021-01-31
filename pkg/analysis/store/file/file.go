package file

// The file implementation of a Store expects of a JSON lines of activities

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fastjson"

	"github.com/bzimmer/gravl/pkg/analysis/store"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

type Option func(*file) error

type file struct {
	path       string
	dirty      bool
	flush      bool
	mutex      sync.RWMutex
	activities map[int64]*strava.Activity
}

// Flush configures whether to flush the contents of any updates (add or remove) to disk on close
func Flush(flush bool) Option {
	return func(f *file) error {
		f.flush = flush
		return nil
	}
}

// Open a file of json-encoded activities
func Open(path string, opts ...Option) (store.Store, error) {
	acts, err := read(path)
	if err != nil {
		return nil, err
	}
	f := &file{path: path, activities: acts}
	for i := range opts {
		if err := opts[i](f); err != nil {
			return nil, err
		}
	}
	return f, nil
}

// Close the file
func (s *file) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.flush && s.dirty {
		// @todo(bzimmer) flush contents to disk
		log.Warn().Msg("contents not flushed")
		return store.UnsupportedOperation
	}
	return nil
}

// Activities returns a slice of (potentially incomplete) Activity instances
func (s *file) Activities(ctx context.Context) (<-chan *strava.Activity, <-chan error) {
	errs := make(chan error)
	acts := make(chan *strava.Activity)
	go func() {
		defer close(acts)
		defer close(errs)
		s.mutex.RLock()
		defer s.mutex.RUnlock()
		for _, act := range s.activities {
			acts <- act
		}
	}()
	return acts, errs
}

// Activity returns a fully populated Activity
func (s *file) Activity(ctx context.Context, activityID int64) (*strava.Activity, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	act, ok := s.activities[activityID]
	if !ok {
		return nil, fmt.Errorf("id not found {%d}", activityID)
	}
	return act, nil
}

// Exists returns true if the activity exists, false otherwise
func (s *file) Exists(ctx context.Context, activityID int64) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	_, ok := s.activities[activityID]
	return ok, nil
}

// Save the activities to the source
func (s *file) Save(ctx context.Context, acts ...*strava.Activity) error {
	if len(acts) == 0 {
		return nil
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, act := range acts {
		s.activities[act.ID] = act
	}
	s.dirty = true
	return nil
}

// Remove the activities from the source
func (s *file) Remove(ctx context.Context, acts ...*strava.Activity) error {
	if len(acts) == 0 {
		return nil
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, act := range acts {
		delete(s.activities, act.ID)
	}
	s.dirty = true
	return nil
}

func read(path string) (map[int64]*strava.Activity, error) {
	var b []byte
	var err error
	var sc fastjson.Scanner
	var activities = make(map[int64]*strava.Activity)
	log.Debug().Str("path", path).Msg("reading")
	b, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	sc.InitBytes(b)
	for sc.Next() {
		if err = sc.Error(); err != nil {
			return nil, err
		}
		val := sc.Value()
		act := &strava.Activity{}
		err = json.Unmarshal(val.MarshalTo(nil), act)
		if err != nil {
			return nil, err
		}
		activities[act.ID] = act
	}
	return activities, nil
}
