package strava

import (
	"context"

	"github.com/bzimmer/gravl/pkg/analysis/store"
	"github.com/bzimmer/gravl/pkg/providers/activity"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

type api struct {
	client *strava.Client
}

// Open a source backed by the Strava API client
func Open(client *strava.Client) store.Source {
	return &api{client: client}
}

// Close is a no-op
func (s *api) Close() error {
	return nil
}

// Activities returns a slice of (potentially incomplete) Activity instances
func (s *api) Activities(ctx context.Context) ([]*strava.Activity, error) {
	return s.client.Activity.Activities(ctx, activity.Pagination{})
}

// Activity returns a fully populated Activity
func (s *api) Activity(ctx context.Context, activityID int64) (*strava.Activity, error) {
	return s.client.Activity.Activity(ctx, activityID)
}
