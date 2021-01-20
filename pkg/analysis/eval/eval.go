package eval

import (
	"context"

	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

type Evaluator interface {

	// Filter the collection of activities by the expression returning those evaluating to true
	Filter(ctx context.Context, acts []*strava.Activity) ([]*strava.Activity, error)

	// Map over the collection of activities producing a slice of expression evaluation values
	Map(ctx context.Context, acts []*strava.Activity) ([]interface{}, error)
}
