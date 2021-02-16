package file

// The file implementation of a Store expects of a JSON lines of activities

import (
	"bufio"
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"

	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
	"github.com/bzimmer/gravl/pkg/store"
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
		fo, err := ioutil.TempFile("", "")
		if err != nil {
			return err
		}
		defer os.Remove(fo.Name())
		writer := bufio.NewWriter(fo)
		enc := json.NewEncoder(writer)
		enc.SetIndent("", " ")
		enc.SetEscapeHTML(false)
		for i := range s.activities {
			if err := enc.Encode(s.activities[i]); err != nil {
				return err
			}
		}
		if err := writer.Flush(); err != nil {
			return err
		}
		if err := fo.Close(); err != nil {
			return err
		}
		if err := os.Rename(fo.Name(), s.path); err != nil {
			return err
		}
		s.dirty = false
		return nil
	}
	return nil
}

// Activities returns a channel of activities and errors for an athlete
func (s *file) Activities(ctx context.Context) <-chan *strava.ActivityResult {
	acts := make(chan *strava.ActivityResult)
	go func() {
		defer close(acts)
		s.mutex.RLock()
		defer s.mutex.RUnlock()
		for _, act := range s.activities {
			select {
			case <-ctx.Done():
				acts <- &strava.ActivityResult{Err: ctx.Err()}
				return
			case acts <- &strava.ActivityResult{Activity: act}:
			}
		}
	}()
	return acts
}

// Activity returns a fully populated Activity
func (s *file) Activity(ctx context.Context, activityID int64) (*strava.Activity, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	act, ok := s.activities[activityID]
	if !ok {
		return nil, store.ErrNotFound
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
	var activities = make(map[int64]*strava.Activity)
	log.Info().Str("path", path).Msg("reading")
	b, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = nil
	gjson.ForEachLine(string(b), func(res gjson.Result) bool {
		act := &strava.Activity{}
		err = json.Unmarshal([]byte(res.Raw), act)
		if err != nil {
			return false
		}
		activities[act.ID] = act
		return true
	})
	if err != nil {
		return nil, err
	}
	log.Info().Str("path", path).Int("activities", len(activities)).Msg("reading")
	return activities, nil
}
