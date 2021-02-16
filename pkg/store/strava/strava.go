package strava

import (
	"context"

	"github.com/bzimmer/gravl/pkg/providers/activity"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
	"github.com/bzimmer/gravl/pkg/store"
)

type api struct {
	client *strava.Client
}

// Open a source backed by the Strava API client
func Open(client *strava.Client) store.Store {
	return &api{client: client}
}

// Close is a no-op
func (s *api) Close() error {
	return nil
}

// Activities returns a channel of activities and errors for an athlete
func (s *api) Activities(ctx context.Context) <-chan *strava.ActivityResult {
	return s.client.Activity.Activities(ctx, activity.Pagination{})
}

// Activity returns a fully populated Activity
func (s *api) Activity(ctx context.Context, activityID int64) (*strava.Activity, error) {
	return s.client.Activity.Activity(ctx, activityID)
}

// Exists returns true if the activity exists, false otherwise
func (s *api) Exists(ctx context.Context, activityID int64) (bool, error) {
	act, err := s.client.Activity.Activity(ctx, activityID)
	if err != nil {
		return false, err
	}
	return act != nil, nil
}

// Save the activities to the source
func (s *api) Save(ctx context.Context, acts ...*strava.Activity) error {
	return store.ErrUnsupportedOperation
}

// Remove the activities from the source
func (s *api) Remove(ctx context.Context, acts ...*strava.Activity) error {
	return store.ErrUnsupportedOperation
}
