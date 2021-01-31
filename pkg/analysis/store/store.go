package store

import (
	"context"
	"errors"

	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

// UnsupportedOperation signals the operation is not supported but this implementation
var UnsupportedOperation = errors.New("unsupported Operation")

type Store interface {

	// Activities returns channels for iteratring through activities
	//
	// The activity returned on the activity channel might not be fully populated per
	//  the implementation details of the specific store. To ensure a fully populated
	//  activity call Activity() with the ID
	//
	// Errors will be sent to the error channel
	//
	// Both channels will be closed on the occurrence of the first error or all activities
	//  have been iterated
	Activities(ctx context.Context) (<-chan *strava.Activity, <-chan error)

	// Activity returns a fully populated Activity
	Activity(ctx context.Context, activityID int64) (*strava.Activity, error)

	// Exists returns true if the activity exists, false otherwise
	Exists(ctx context.Context, activityID int64) (bool, error)

	// Save the activities to the source
	Save(ctx context.Context, acts ...*strava.Activity) error

	// Remove the activities from the source
	Remove(ctx context.Context, acts ...*strava.Activity) error

	// Close the source
	Close() error
}
