package analysis

import (
	"context"
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/spf13/cast"

	"github.com/bzimmer/gravl/pkg/activity/strava"
)

// Pass represents a collection of Activities for analysis
type Pass struct {
	// Units of the resulting Activities
	Units Units
	// Activities on which analysis will occur
	Activities []*strava.Activity
}

// Group represents a single group in a group tree
type Group struct {
	// Key is the result of apply an expression against an Activity
	Key string
	// Pass holds the Activities grouped by Key
	Pass *Pass
	// Groups holds child Groups if more than one level of grouping exists
	Groups []*Group
	// Level of the group in the tree
	Level int
}

func (g *Group) walk(ctx context.Context, f func(context.Context, *Group) error) error {
	if err := f(ctx, g); err != nil {
		return err
	}
	for i := range g.Groups {
		if err := g.Groups[i].walk(ctx, f); err != nil {
			return err
		}
	}
	return nil
}

// Filter filters the activities using an expression
// For example:
//  {.Type in ["Ride"] && !.Commute && .StartDateLocal.Year() in [2020, 2019]}
func (p *Pass) Filter(q string) (*Pass, error) {
	code := fmt.Sprintf("filter(Activities, %s)", q)
	out, err := expr.Eval(code, p)
	if err != nil {
		return nil, err
	}
	res := out.([]interface{})
	acts := make([]*strava.Activity, len(res))
	for i := range res {
		acts[i] = res[i].(*strava.Activity)
	}
	return &Pass{Activities: acts, Units: p.Units}, nil
}

// GroupBy groups activities by a key
func (p *Pass) GroupBy(exprs ...string) (*Group, error) {
	g := &Group{Pass: p}
	if err := groupby(g, exprs...); err != nil {
		return nil, err
	}
	return g, nil
}

func groupby(group *Group, exprs ...string) error {
	if len(exprs) == 0 {
		return nil
	}
	// map over the activities to generate a group key
	q := exprs[0]
	code := fmt.Sprintf("map(Activities, %s)", q)
	out, err := expr.Eval(code, group.Pass)
	if err != nil {
		return err
	}
	// group all activities into a Group based on their group key
	res := out.([]interface{})
	passes := make(map[string]*Pass, len(res))
	for i, k := range res {
		var key string
		key, err = cast.ToStringE(k)
		if err != nil {
			return err
		}
		if _, ok := passes[key]; !ok {
			passes[key] = &Pass{Units: group.Pass.Units}
		}
		passes[key].Activities = append(passes[key].Activities, group.Pass.Activities[i])
	}
	// recurse if more grouping operators exist
	tail := exprs[1:]
	for key, pass := range passes {
		parent := &Group{Key: key, Pass: pass, Level: group.Level + 1}
		group.Groups = append(group.Groups, parent)
		if len(tail) > 0 {
			if err = groupby(parent, tail...); err != nil {
				return err
			}
		}
	}
	return nil
}
