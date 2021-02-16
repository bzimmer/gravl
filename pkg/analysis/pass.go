package analysis

import (
	"context"
	"fmt"

	"github.com/bzimmer/gravl/pkg/eval"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

// Pass represents a collection of Activities for analysis
type Pass struct {
	// Key is the result of applying an expression on an Activity
	Key string `json:"key"`
	// Activities on which analysis will occur
	Activities []*strava.Activity `json:"activities"`
	// Children contains all the child Passes
	Children []*Pass `json:"children"`
}

// Group groups activities by result of evaluating the list of mappers
func Group(ctx context.Context, acts []*strava.Activity, mappers ...eval.Mapper) (*Pass, error) {
	pass := &Pass{Activities: acts}
	return pass, group(ctx, pass, mappers)
}

func group(ctx context.Context, pass *Pass, evals []eval.Mapper) error {
	if len(evals) == 0 {
		return nil
	}
	keys, err := evals[0].Map(ctx, pass.Activities)
	if err != nil {
		return err
	}
	groups := make(map[string][]*strava.Activity)
	for i := range keys {
		key := fmt.Sprintf("%v", keys[i])
		groups[key] = append(groups[key], pass.Activities[i])
	}
	for key, acts := range groups {
		child := &Pass{Activities: acts, Key: key}
		if err := group(ctx, child, evals[1:]); err != nil {
			return err
		}
		pass.Children = append(pass.Children, child)
	}
	return nil
}
