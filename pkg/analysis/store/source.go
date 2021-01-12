package store

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/valyala/fastjson"

	"github.com/bzimmer/gravl/pkg/providers/activity"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

type Source interface {
	// Activities returns a slice of (potentially incomplete) Activity instances
	Activities(ctx context.Context) ([]*strava.Activity, error)

	// Activity returns a fully populated Activity
	Activity(ctx context.Context, activityID int64) (*strava.Activity, error)
}

type SourceFile struct {
	Path       string
	activities map[int64]*strava.Activity
}

func (s *SourceFile) Activities(ctx context.Context) ([]*strava.Activity, error) {
	var b []byte
	var err error
	var sc fastjson.Scanner

	s.activities = make(map[int64]*strava.Activity)

	b, err = ioutil.ReadFile(s.Path)
	if err != nil {
		return nil, err
	}
	sc.InitBytes(b)
	for sc.Next() {
		if err = sc.Error(); err != nil {
			return nil, err
		}
		val := sc.Value()
		act := &strava.Activity{}
		err = json.Unmarshal(val.MarshalTo(nil), act)
		if err != nil {
			return nil, err
		}
		s.activities[act.ID] = act
	}

	var i int
	acts := make([]*strava.Activity, len(s.activities))
	for _, act := range s.activities {
		acts[i] = act
		i++
	}
	return acts, nil
}

func (s *SourceFile) Activity(ctx context.Context, activityID int64) (*strava.Activity, error) {
	act, ok := s.activities[activityID]
	if !ok {
		return nil, fmt.Errorf("id not found {%d}", activityID)
	}
	return act, nil
}

type SourceStrava struct {
	Client *strava.Client
}

func (s *SourceStrava) Activities(ctx context.Context) ([]*strava.Activity, error) {
	return s.Client.Activity.Activities(ctx, activity.Pagination{})
}

func (s *SourceStrava) Activity(ctx context.Context, activityID int64) (*strava.Activity, error) {
	return s.Client.Activity.Activity(ctx, activityID)
}
