package memory

import (
	"context"
	"sync"

	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
	"github.com/bzimmer/gravl/pkg/store"
)

// Provider is a generator of strava activities
type Provider interface {
	// Activities generates strava activities
	Activities() ([]*strava.Activity, error)
	// Close is called when the store is closed
	Close(map[int64]*strava.Activity) error
}

type memory struct {
	dirty      bool
	mutex      sync.RWMutex
	provider   Provider
	activities map[int64]*strava.Activity
}

// Open the store supplied data provivder
func Open(provider Provider) (store.Store, error) {
	acts, err := provider.Activities()
	if err != nil {
		return nil, err
	}
	m := &memory{activities: make(map[int64]*strava.Activity), provider: provider}
	for i := range acts {
		m.activities[acts[i].ID] = acts[i]
	}
	return m, nil
}

// Close the store
func (s *memory) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.dirty {
		if err := s.provider.Close(s.activities); err != nil {
			return err
		}
	}
	s.dirty = false
	return nil
}

// Activities returns a channel of activities and errors for an athlete
func (s *memory) Activities(ctx context.Context) <-chan *strava.ActivityResult {
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
func (s *memory) Activity(ctx context.Context, activityID int64) (*strava.Activity, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	act, ok := s.activities[activityID]
	if !ok {
		return nil, store.ErrNotFound
	}
	return act, nil
}

// Exists returns true if the activity exists, false otherwise
func (s *memory) Exists(ctx context.Context, activityID int64) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	_, ok := s.activities[activityID]
	return ok, nil
}

// Save the activities to the source
func (s *memory) Save(ctx context.Context, acts ...*strava.Activity) error {
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
func (s *memory) Remove(ctx context.Context, acts ...*strava.Activity) error {
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
