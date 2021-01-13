package store

import (
	"context"

	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

type Closer interface {

	// Close the source
	Close() error
}

type Source interface {
	Closer

	// Activities returns a slice of (potentially incomplete) Activity instances
	Activities(ctx context.Context) ([]*strava.Activity, error)

	// Activity returns a fully populated Activity
	Activity(ctx context.Context, activityID int64) (*strava.Activity, error)
}

type Sink interface {
	Closer

	// Exists returns true if the activity exists, false otherwise
	Exists(ctx context.Context, activityID int64) (bool, error)

	// Save the activities to the source
	Save(ctx context.Context, acts ...*strava.Activity) error

	// Remove the activities from the source
	Remove(ctx context.Context, acts ...*strava.Activity) error
}

type SourceSink interface {
	Source
	Sink
}
