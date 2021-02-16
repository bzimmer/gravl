package store

import (
	"context"
	"errors"

	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

// ErrUnsupportedOperation is returned when the operation is not supported
var ErrUnsupportedOperation = errors.New("unsupported operation")

// ErrNotFound is returned when the activity is not found
var ErrNotFound = errors.New("activity is not found")

type Store interface {

	// Activities returns a channel of activities and errors for an athlete
	//
	// The activity returned on the channel might not be fully populated per the
	//  implementation details of the specific store. To ensure a fully populated
	//  activity call `Activity()` with the ID
	Activities(ctx context.Context) <-chan *strava.ActivityResult

	// Activity returns a fully populated Activity
	Activity(ctx context.Context, activityID int64) (*strava.Activity, error)

	// Exists returns true if the activity exists, false otherwise
	Exists(ctx context.Context, activityID int64) (bool, error)

	// Save the activities to the store
	Save(ctx context.Context, acts ...*strava.Activity) error

	// Remove the activities from the store
	Remove(ctx context.Context, acts ...*strava.Activity) error

	// Close the store
	Close() error
}
