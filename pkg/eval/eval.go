package eval

import (
	"context"

	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

// Filterer performs activity filtering
type Filterer interface {
	// Filter the collection of activities by the expression returning those evaluating to true
	Filter(ctx context.Context, acts []*strava.Activity) ([]*strava.Activity, error)
}

// Mapper performs activity mapping
type Mapper interface {
	// Map over the collection of activities producing a slice of expression evaluation values
	Map(ctx context.Context, acts []*strava.Activity) ([]interface{}, error)
}

// Evalutor performs evaluations on activities
type Evaluator interface {
	// Bool performs an evaluation resulting in a boolean value
	Bool(ctx context.Context, act *strava.Activity) (bool, error)

	// Eval performs an evaluation on an activity with an arbitrary result
	Eval(ctx context.Context, act *strava.Activity) (interface{}, error)
}
