package eval

import (
	"context"

	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

type Evaluator interface {

	// Filter the collection of activities by the expression returning those evaluating to true
	Filter(ctx context.Context, q string, acts []*strava.Activity) ([]*strava.Activity, error)

	// GroupBy the collection of activities by the result of applying the expression
	GroupBy(ctx context.Context, q string, acts []*strava.Activity) (map[string][]*strava.Activity, error)
}
